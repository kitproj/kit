# Dynamo

Dynamo is a distributed key-value storage system developed by Amazon Web Services (AWS) for highly available, scalable and reliable NoSQL data storage.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: amazon/dynamodb-local
    name: dynamo
    ports: "8000"
```

