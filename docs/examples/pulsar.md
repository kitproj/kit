# Pulsar

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - command: /pulsar/bin/pulsar standalone
    image: apachepulsar/pulsar
    name: pulsar
    ports: 6650 8080
```

