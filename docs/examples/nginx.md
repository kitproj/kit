# Nginx

> Maintainer: NGINX Docker Maintainers <docker-maint@nginx.com>

```yaml
apiVersion: kit/v1
kind: Tasks
metadata:
  name: example
spec:
  tasks:
  - image: nginx
    name: nginx
    ports: 80:8080
    volumeMounts:
    - mountPath: /usr/share/nginx/html
      name: nginx.html
  volumes:
  - hostPath:
      path: volumes/nginx/html
    name: nginx.html
```

