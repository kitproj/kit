# Rabbitmq

The RabbitMQ container image provides a scalable and highly-available message broker that supports various messaging protocols and can be easily deployed and managed in a containerized environment.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: rabbitmq
    name: rabbitmq
    ports: 4369 5671 5672 15691 15692 25672
    volumeMounts:
    - mountPath: /var/lib/rabbitmq
      name: rabbitmq.rabbitmq
  volumes:
  - hostPath:
      path: volumes/rabbitmq/rabbitmq
    name: rabbitmq.rabbitmq
```

