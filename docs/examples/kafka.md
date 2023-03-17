# Kafka  with Kraft

[Help](https://github.com/kitproj/kafka-docker)

Standalone Kafka image that uses Kraft.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: ghcr.io/kitproj/kafka
    name: kafka
    ports: "9092"
    volumeMounts:
    - mountPath: /tmp/kraft-combined-logs
      name: kafka.kraft-combined-logs
  volumes:
  - hostPath:
      path: volumes/kafka/kraft-combined-logs
    name: kafka.kraft-combined-logs
```

Licence(s): MIT

