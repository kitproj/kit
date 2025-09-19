# Sim

[Help](https://github.com/kitproj/sim)

Sim is straight-forward API simulation tool that's tiny, fast, secure and scalable.

```yaml
tasks:
  "":
    image: ghcr.io/kitproj/sim
    ports:
    - "8080"
    readinessProbe: http://localhost:8080/hello
volumes:
- hostPath:
    path: volumes/sim/apis
  name: sim.apis
```

Licence(s): MIT

