# Rabbitmq

The RabbitMQ container image provides a scalable and highly-available message broker that supports various messaging protocols and can be easily deployed and managed in a containerized environment.

```yaml
tasks:
  "":
    image: rabbitmq
volumes:
- hostPath:
    path: volumes/rabbitmq/rabbitmq
  name: rabbitmq.rabbitmq
```

