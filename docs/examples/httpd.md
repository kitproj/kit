# Httpd

The httpd container image provides an Apache HTTP server for serving web content.

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: httpd
    name: httpd
    ports: 80:8080
    volumeMounts:
    - mountPath: /usr/local/apache2/htdocs
      name: apache2.htdocs
  volumes:
  - hostPath:
      path: volumes/httpd/htdocs
    name: httpd.apis
```

