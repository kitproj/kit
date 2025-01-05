# Nginx

The nginx container image is a lightweight and high-performance web server designed to efficiently serve static and dynamic content.

> Maintainer: NGINX Docker Maintainers <docker-maint@nginx.com>

```yaml
tasks:
  "":
    image: nginx
    volumeMounts:
    - mountPath: /usr/share/nginx/html
      name: nginx.html
volumes:
- hostPath:
    path: volumes/nginx/html
  name: nginx.html
```

