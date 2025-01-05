# Mongo

The mongo container image provides a lightweight and scalable solution for running MongoDB databases in a containerized environment.

```yaml
tasks:
  "":
    image: mongo
volumes:
- hostPath:
    path: volumes/mongo/configdb
  name: mongo.configdb
- hostPath:
    path: volumes/mongo/db
  name: mongo.db
```

