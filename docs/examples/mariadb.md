# MariaDB Database

[Help](https://github.com/MariaDB/mariadb-docker)

MariaDB Database for relational SQL

https://hub.docker.com/_/mariadb/

> Maintainer: MariaDB Community

```yaml
tasks:
  "":
    env:
    - MARIADB_ROOT_PASSWORD=password
    image: mariadb
volumes:
- hostPath:
    path: volumes/mariadb/mysql
  name: mariadb.mysql
```

Licence(s): GPL-2.0

