
#### Container Runtimes (CRI)

```yaml
cri:
 version: 18.6.0
 type: docker
 config:
   log-driver: json-file
   log-opts:
     max-size: 1000m
     max-file": 3
```
