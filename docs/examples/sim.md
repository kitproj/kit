# Sim

[Help](https://github.com/kitproj/sim)

Sim is straight-forward API simulation tool that's tiny, fast, secure and scalable.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: ghcr.io/kitproj/sim
    name: sim
    ports: "8080"
    readinessProbe: http://:8080/hello?failureThreshold=20&initialDelay=3s&period=3s&successThreshold=1
    volumeMounts:
    - mountPath: /apis
      name: sim.apis
  volumes:
  - hostPath:
      path: volumes/sim/apis
    name: sim.apis
```

Licence(s): MIT

