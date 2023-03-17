# Mongo

The mongo container image provides a lightweight and scalable solution for running MongoDB databases in a containerized environment.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: mongo
    name: mongo
    ports: "27017"
    volumeMounts:
    - mountPath: /data/configdb
      name: mongo.configdb
    - mountPath: /data/db
      name: mongo.db
  volumes:
  - hostPath:
      path: volumes/mongo/configdb
    name: mongo.configdb
  - hostPath:
      path: volumes/mongo/db
    name: mongo.db
```

