# Mysql

The MySQL container image provides a pre-configured and optimized environment for running the MySQL database service within a container.

```yaml
tasks:
  "":
    env:
      MYSQL_ROOT_PASSWORD: password
    image: mysql
volumes:
- hostPath:
    path: volumes/mysql/mysql
  name: mysql.mysql
```

