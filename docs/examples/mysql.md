# Mysql

The MySQL container image provides a pre-configured and optimized environment for running the MySQL database service within a container.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - env:
    - MYSQL_ROOT_PASSWORD=password
    image: mysql
    name: mysql
    ports: 3306 33060
    volumeMounts:
    - mountPath: /var/lib/mysql
      name: mysql.mysql
  volumes:
  - hostPath:
      path: volumes/mysql/mysql
    name: mysql.mysql
```

