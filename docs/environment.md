### Environment
Environment variables are saved to `/etc/environment/` and are sourced before any commands runs.
```yaml
environment:
  env1: "val: {{env1}}"
```
