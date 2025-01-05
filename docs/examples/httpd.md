# Httpd

The httpd container image provides an Apache HTTP server for serving web content.

```yaml
tasks:
  "":
    image: httpd
    volumeMounts:
    - mountPath: /usr/local/apache2/htdocs
      name: apache2.htdocs
volumes:
- hostPath:
    path: volumes/httpd/htdocs
  name: httpd.apis
```

