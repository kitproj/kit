package proc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/adler32"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/kitproj/kit/internal/metrics"
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
	"k8s.io/utils/strings/slices"
	"sigs.k8s.io/yaml"
)

type k8s struct {
	log  *log.Logger
	spec types.Spec
	name string
	pods []string // namespace/name
	types.Task
}

// previously we used the K8s common labels, but because charts use them themselves (e.g. Helm) we cannot and must create our own annotations
const x = "kit.kitproj.github.com"
const nameLabel = x + "/name"
const versionLabel = x + "/version"

func (k *k8s) Run(ctx context.Context, stdout io.Writer, stderr io.Writer) error {

	log := k.log
	// apply the manifests
	var files []string
	for _, manifest := range k.Manifests {
		file := filepath.Join(k.WorkingDir, manifest)
		// walk the file tree
		err := filepath.WalkDir(file, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("failed to walk directory: %w", err)
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
			return fmt.Errorf("failed to walk directory: %w", err)
		}
	}

	// connect to the k8s cluster
	kubeConfig := os.Getenv("KUBECONFIG")
	if kubeConfig == "" {
		kubeConfig = clientcmd.RecommendedHomeFile
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	// Get the namespace associated with the current context
	defaultNamespace, _, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig},
		&clientcmd.ConfigOverrides{},
	).Namespace()
	if err != nil {
		return fmt.Errorf("failed to get namespace: %w", err)
	}

	if k.Namespace != "" {
		defaultNamespace = k.Namespace
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	// Create a Discovery client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create discovery client: %w", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// for each manifest, read it as YAML (splitting by ---)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		var uns []*unstructured.Unstructured

		// if the YAML contains multiple documents, split them
		for _, doc := range bytes.Split(data, []byte("\n---\n")) {
			var manifest map[string]any
			err = yaml.Unmarshal(doc, &manifest)
			if err != nil {
				return fmt.Errorf("failed to unmarshal YAML: %w", err)
			}
			uns = append(uns, &unstructured.Unstructured{Object: manifest})
		}

		sortUnstructureds(uns)

		// for each YAML document, create the object
		for _, u := range uns {
			apiResources, err := discoveryClient.ServerResourcesForGroupVersion(u.GetAPIVersion())
			if err != nil {
				return fmt.Errorf("failed to get server resources: %w", err)
			}

			// Find the resource that matches the kind
			var resource string
			var namespaced bool
			kind := u.GetKind()
			for _, apiResource := range apiResources.APIResources {
				if apiResource.Kind == kind {
					resource = apiResource.Name
					namespaced = apiResource.Namespaced
					break
				}
			}

			gvr := schema.GroupVersionResource{
				Group:    u.GroupVersionKind().Group,
				Version:  u.GroupVersionKind().Version,
				Resource: resource,
			}

			if u.GetLabels() == nil {
				u.SetLabels(make(map[string]string))
			}
			if u.GetAnnotations() == nil {
				u.SetAnnotations(make(map[string]string))
			}
			labels := u.GetLabels()
			labels[nameLabel] = k.name
			u.SetLabels(labels)

			// if this is a deployment or a statefulset, then add the label to the pod template
			if u.GetKind() == "Deployment" || u.GetKind() == "StatefulSet" {
				// update selector labels
				labels, _, err := unstructured.NestedMap(u.Object, "spec", "selector", "matchLabels")
				if err != nil {
					return fmt.Errorf("failed to get selector labels: %w", err)
				}
				labels[nameLabel] = k.name
				err = unstructured.SetNestedMap(u.Object, labels, "spec", "selector", "matchLabels")
				if err != nil {
					return fmt.Errorf("failed to set selector labels: %w", err)
				}

				// update template labels
				labels, _, err = unstructured.NestedMap(u.Object, "spec", "template", "metadata", "labels")
				if err != nil {
					return fmt.Errorf("failed to get template labels: %w", err)
				}
				labels[nameLabel] = k.name
				err = unstructured.SetNestedMap(u.Object, labels, "spec", "template", "metadata", "labels")
				if err != nil {
					return fmt.Errorf("failed to set template labels: %w", err)
				}
			}

			if namespaced && u.GetNamespace() == "" {
				u.SetNamespace(defaultNamespace)
			}

			// add a hash of the manifest to the annotations
			annotations := u.GetAnnotations()
			data, _ := json.Marshal(u.Object)
			annotations[versionLabel] = fmt.Sprintf("%x", adler32.Checksum(data))
			u.SetAnnotations(annotations)

			// has it been created already?
			existing, err := dynamicClient.Resource(gvr).Namespace(u.GetNamespace()).Get(ctx, u.GetName(), metav1.GetOptions{})
			if err != nil {
				if !apierrors.IsNotFound(err) {
					return fmt.Errorf("failed to get resource: %w", err)
				}
			} else {
				expectedHash := u.GetAnnotations()[versionLabel]
				// has the manifest changed?
				existingHash := existing.GetAnnotations()[versionLabel]
				if existingHash == expectedHash {
					log.Printf("%s/%s/%s unchanged\n", resource, u.GetNamespace(), u.GetName())
					continue
				}

				log.Printf("deleting %s/%s/%s\n", resource, u.GetNamespace(), u.GetName())

				err = dynamicClient.Resource(gvr).Namespace(u.GetNamespace()).Delete(ctx, u.GetName(), metav1.DeleteOptions{})
				if err != nil {
					return fmt.Errorf("failed to delete resource: %w", err)
				}
				// wait for the resource to be deleted
				for {
					_, err = dynamicClient.Resource(gvr).Namespace(u.GetNamespace()).Get(ctx, u.GetName(), metav1.GetOptions{})
					if apierrors.IsNotFound(err) {
						break
					}
					time.Sleep(1 * time.Second)
				}
			}

			log.Printf("creating %s/%s/%s\n", resource, u.GetNamespace(), u.GetName())

			_, err = dynamicClient.Resource(gvr).Namespace(u.GetNamespace()).Create(ctx, u, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create resource: %w", err)
			}
		}
	}

	ports := k.Ports.Map()

	// we can exit if we are not expecting to forward any ports
	if len(ports) == 0 {
		return nil
	}

	// Create a shared informer factory for only the labelled resource managed-by kit and named after the task
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 10*time.Second, informers.WithTweakListOptions(func(options *metav1.ListOptions) {
		options.LabelSelector = fmt.Sprintf("%s=%s", nameLabel, k.name)
	}))

	// Create a pod informer
	podInformer := factory.Core().V1().Pods().Informer()

	logging := sync.Map{}        // namespace/name/container -> true
	portForwarding := sync.Map{} // port -> true

	// Add event handlers
	processPod := func(obj any) {
		pod := obj.(*corev1.Pod)

		podKey := pod.Namespace + "/" + pod.Name

		if !slices.Contains(k.pods, podKey) {
			k.pods = append(k.pods, podKey)
		}

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

				// check if the pod is already being logged
				if _, ok := logging.Load(key); ok {
					return
				}

				logging.Store(key, true)
				defer logging.Delete(key)

				defer func() {
					if r := recover(); r != nil {
						log.Printf("error while tailing logs: %s: %v\n", key, r)
					}
				}()

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
				hostPort := ports[uint16(containerPort)]

				if hostPort == 0 {
					continue
				}

				// start port-forwarding
				go func() {
					// check if the pod is already being port-forwarded
					if _, ok := portForwarding.Load(hostPort); ok {
						return
					}

					portForwarding.Store(hostPort, true)
					defer portForwarding.Delete(hostPort)

					defer func() {
						if r := recover(); r != nil {
							log.Printf("error while port-forwarding: %d: %v\n", hostPort, r)
						}
					}()

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

					fw, err := portforward.New(dialer, ports, stopChan, readyChan, nil, nil)
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
		return fmt.Errorf("failed to add event handler: %w", err)
	}

	factory.Start(ctx.Done())

	<-ctx.Done()

	return nil

}

func sortUnstructureds(uns []*unstructured.Unstructured) {
	// we need to sort the unstructured outputs by their kind, so that namespaces get applied before deployments, etc
	// much like Helm/Argo CD does
	// this is because some resources depend on others, e.g. a deployment depends on a namespace
	order := []string{
		"APIService",
		"Ingress",
		"Service",
		"CronJob",
		"Job",
		"StatefulSet",
		"Deployment",
		"ReplicaSet",
		"ReplicationController",
		"Pod",
		"DaemonSet",
		"RoleBinding",
		"Role",
		"ClusterRoleBinding",
		"ClusterRole",
		"CustomResourceDefinition",
		"ServiceAccount",
		"PersistentVolumeClaim",
		"PersistentVolume",
		"StorageClass",
		"ConfigMap",
		"Secret",
		"PodSecurityPolicy",
		"LimitRange",
		"ResourceQuota",
		"Namespace",
	}
	sort.Slice(uns, func(i, j int) bool {
		// slices.Index will return -1 if the element is not found
		return slices.Index(order, uns[i].GetKind()) > slices.Index(order, uns[j].GetKind())
	})
}

func (k *k8s) GetMetrics(ctx context.Context) (*types.Metrics, error) {
	sum := &types.Metrics{}
	for _, podKey := range k.pods {
		parts := strings.SplitN(podKey, "/", 2)
		namespace := parts[0]
		podName := parts[1]
		metrics, err := k.getMetrics(ctx, namespace, podName)
		if err != nil {
			return nil, err
		}
		sum.CPU += metrics.CPU
		sum.Mem += metrics.Mem
	}
	return sum, nil
}

func (k *k8s) getMetrics(ctx context.Context, namespace, podName string) (*types.Metrics, error) {
	command, args := metrics.GetAllProcessesCommand()
	cmdArgs := append([]string{"exec", "-n", namespace, podName, "--"}, command)
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.CommandContext(ctx, "kubectl", cmdArgs...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("kubectl exec ps failed for pod %s/%s: %w", namespace, podName, err)
	}

	return metrics.ParseOutput(string(output))
}


