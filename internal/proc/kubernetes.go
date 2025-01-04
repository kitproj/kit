package proc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash/adler32"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kitproj/kit/internal/types"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"sigs.k8s.io/yaml"
)

type k8s struct {
	log  *log.Logger
	spec types.PodSpec
	types.Task
}

func (k *k8s) Run(ctx context.Context, stdout io.Writer, stderr io.Writer) error {

	log := k.log
	// apply the manifests
	var files []string
	for _, manifest := range k.Manifests {
		file := filepath.Join(k.WorkingDir, manifest)
		// walk the file tree
		err := filepath.WalkDir(file, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// must be a YAML ile
			if d.IsDir() {
				return nil
			}
			ext := filepath.Ext(path)
			if ext != ".yaml" && ext != ".yml" {
				return nil
			}
			files = append(files, path)
			return nil
		})
		if err != nil {
			return err
		}
	}

	// connect to the k8s cluster
	kubeConfig := os.Getenv("KUBECONFIG")
	if kubeConfig == "" {
		kubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return err
	}

	// Get the namespace associated with the current context
	defaultNamespace, _, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: clientcmd.RecommendedHomeFile},
		&clientcmd.ConfigOverrides{},
	).Namespace()
	if err != nil {
		return err
	}

	log.Printf("Using namespace %s\n", defaultNamespace)

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// Create a Discovery client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return err
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	// for each manifest, read it as YAML (splitting by ---)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		var uns []*unstructured.Unstructured

		// if the YAML contains multiple documents, split them
		for _, doc := range bytes.Split(data, []byte("\n---\n")) {
			var manifest map[string]any
			err = yaml.Unmarshal(doc, &manifest)
			if err != nil {
				return err
			}
			u := &unstructured.Unstructured{Object: manifest}
			if u.GetLabels() == nil {
				u.SetLabels(make(map[string]string))
			}
			if u.GetAnnotations() == nil {
				u.SetAnnotations(make(map[string]string))
			}
			labels := u.GetLabels()
			labels["app.kubernetes.io/managed-by"] = "kit"
			labels["app.kubernetes.io/name"] = k.Name
			u.SetLabels(labels)

			if u.GetNamespace() == "" {
				u.SetNamespace(defaultNamespace)
			}

			// if this is a deployment or a statefulset, then add the label to the pod template
			if u.GetKind() == "Deployment" || u.GetKind() == "StatefulSet" {
				// update selector labels
				labels, _, err := unstructured.NestedMap(u.Object, "spec", "selector", "matchLabels")
				if err != nil {
					return err
				}
				labels["app.kubernetes.io/managed-by"] = "kit"
				labels["app.kubernetes.io/name"] = k.Name
				err = unstructured.SetNestedMap(u.Object, labels, "spec", "selector", "matchLabels")
				if err != nil {
					return err
				}

				// update template labels
				labels, _, err = unstructured.NestedMap(u.Object, "spec", "template", "metadata", "labels")
				if err != nil {
					return err
				}
				labels["app.kubernetes.io/managed-by"] = "kit"
				labels["app.kubernetes.io/name"] = k.Name
				err = unstructured.SetNestedMap(u.Object, labels, "spec", "template", "metadata", "labels")
				if err != nil {
					return err
				}
			}

			// add a hash of the manifest to the annotations
			annotations := u.GetAnnotations()
			annotations["app.kubernetes.io/version"] = fmt.Sprintf("%x", adler32.Checksum(doc))
			u.SetAnnotations(annotations)

			uns = append(uns, u)
		}

		// for each YAML document, create the object
		for _, u := range uns {

			apiResources, err := discoveryClient.ServerResourcesForGroupVersion(u.GetAPIVersion())
			if err != nil {
				return err
			}

			// Find the resource that matches the kind
			var resource string
			for _, apiResource := range apiResources.APIResources {
				if apiResource.Kind == u.GetKind() {
					resource = apiResource.Name
					break
				}
			}

			gvr := schema.GroupVersionResource{
				Group:    u.GroupVersionKind().Group,
				Version:  u.GroupVersionKind().Version,
				Resource: resource,
			}

			// has it been created already?
			existing, err := dynamicClient.Resource(gvr).Namespace(u.GetNamespace()).Get(ctx, u.GetName(), metav1.GetOptions{})
			if err != nil {
				if !apierrors.IsNotFound(err) {
					return err
				}
			} else {
				expectedHash := u.GetAnnotations()["app.kubernetes.io/version"]
				// has the manifest changed?
				existingHash := existing.GetAnnotations()["app.kubernetes.io/version"]
				if existingHash == expectedHash {
					log.Printf("Skipping %s/%s/%s: unchanged\n", u.GetAPIVersion(), u.GetKind(), u.GetName())
					continue
				}

				// delete the object
				log.Printf("Deleting %s/%s/%s\n", u.GetAPIVersion(), u.GetKind(), u.GetName())

				err = dynamicClient.Resource(gvr).Namespace(u.GetNamespace()).Delete(ctx, u.GetName(), metav1.DeleteOptions{})
				if err != nil {
					return err
				}
			}

			log.Printf("Creating %s/%s/%s\n", u.GetAPIVersion(), u.GetKind(), u.GetName())

			_, err = dynamicClient.Resource(gvr).Namespace(u.GetNamespace()).Create(ctx, u, metav1.CreateOptions{})
			if err != nil {
				return err
			}
		}
	}

	// Create a shared informer factory
	factory := informers.NewSharedInformerFactory(clientset, 0)

	// Create a pod informer
	podInformer := factory.Core().V1().Pods().Informer()

	logging := sync.Map{}
	portForwarding := sync.Map{}

	// Add event handlers
	processPod := func(obj any) {
		pod := obj.(*corev1.Pod)

		// is the pod labelled with the managed-by label?
		if pod.GetLabels()["app.kubernetes.io/managed-by"] != "kit" {
			return
		}
		if pod.GetLabels()["app.kubernetes.io/name"] != k.Name {
			return
		}

		log.Printf("Pod added: %s/%s\n", pod.Namespace, pod.Name)

		running := make(map[string]bool)

		for _, s := range append(pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses...) {
			running[s.Name] = s.State.Running != nil
		}

		for _, c := range append(pod.Spec.InitContainers, pod.Spec.Containers...) {
			// skip containers that are not running
			if !running[c.Name] {
				continue
			}
			go func() {
				// start a log tail
				key := pod.Namespace + "/" + pod.Name + "/" + c.Name
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Error while tailing logs: %s: %v\n", key, r)
					}
					logging.Delete(key)
				}()

				// check if the pod is already being logged
				if _, ok := logging.Load(key); ok {
					log.Printf("Skipping log tail for container %s: already tailing\n", key)
					return
				}
				logging.Store(key, true)

				log.Printf("Starting log tail for container %s\n", key)

				req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
					Follow:    true,
					Container: c.Name,
					SinceTime: &metav1.Time{Time: time.Now()},
				})
				podLogs, err := req.Stream(ctx)
				if err != nil {
					panic(fmt.Errorf("Error opening stream: %s\n", err))
				}
				defer podLogs.Close()
				_, err = io.Copy(stdout, podLogs)
				if err != nil && !errors.Is(err, context.Canceled) {
					panic(fmt.Errorf("Error copying stream: %s\n", err))
				}
			}()
			for _, port := range c.Ports {
				// only forward host ports
				containerPort := port.ContainerPort
				hostPort := 0

				for _, p := range k.Ports {
					if p.ContainerPort == uint16(containerPort) {
						hostPort = int(p.HostPort)
						break
					}
				}

				if hostPort == 0 {
					log.Printf("Skipping port-forward for container %s/%s port %d\n", pod.Name, c.Name, containerPort)
					continue
				}

				// start port-forwarding
				go func() {
					key := pod.Namespace + "/" + pod.Name + "/" + c.Name + " / " + fmt.Sprintf("%d", containerPort)
					defer func() {
						if r := recover(); r != nil {
							log.Printf("Error while port-forwarding: %s: %v\n", key, r)
						}
						portForwarding.Delete(key)
					}()

					// check if the pod is already being port-forwarded
					if _, ok := portForwarding.Load(key); ok {
						log.Printf("Skipping port-forward for pod %s/%s container %s: already forwarding\n", pod.Namespace, pod.Name, c.Name)
						return
					}

					portForwarding.Store(key, true)

					log.Printf("Starting port-forward for pod %s/%s container %s port %d to host %d\n", pod.Namespace, pod.Name, c.Name, containerPort, hostPort)

					req := clientset.CoreV1().RESTClient().Post().
						Resource("pods").
						Namespace(pod.Namespace).
						Name(pod.Name).
						SubResource("portforward")

					transport, upgrader, err := spdy.RoundTripperFor(config)
					if err != nil {
						panic(err)
					}

					dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", req.URL())

					stopChan := ctx.Done()
					readyChan := make(chan struct{})

					ports := []string{fmt.Sprintf("%d:%d", hostPort, containerPort)}

					fw, err := portforward.New(dialer, ports, stopChan, readyChan, stdout, stderr)
					if err != nil {
						panic(err)
					}

					if err := fw.ForwardPorts(); err != nil {
						panic(err)
					}
				}()
			}
		}
	}
	_, err = podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: processPod,
		UpdateFunc: func(_, newObj any) {
			processPod(newObj)
		},
	})
	if err != nil {
		return err
	}

	factory.Start(ctx.Done())

	<-ctx.Done()

	return nil

}

func (k *k8s) Reset(ctx context.Context) error {
	return nil
}
