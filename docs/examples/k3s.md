# K3s

[Help](https://github.com/k3s-io/k3s)

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: rancher/k3s
    name: k3s
    volumeMounts:
    - mountPath: /var/lib/cni
      name: k3s.cni
    - mountPath: /var/lib/kubelet
      name: k3s.kubelet
    - mountPath: /var/lib/rancher/k3s
      name: k3s.k3s
    - mountPath: /var/log
      name: k3s.log
  volumes:
  - hostPath:
      path: volumes/k3s/cni
    name: k3s.cni
  - hostPath:
      path: volumes/k3s/kubelet
    name: k3s.kubelet
  - hostPath:
      path: volumes/k3s/k3s
    name: k3s.k3s
  - hostPath:
      path: volumes/k3s/log
    name: k3s.log
```

