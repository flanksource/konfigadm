### Containers

```yaml
containers:
  - image: "docker.io/consul:1.3.1"
    args: agent -ui -bootstrap -server
    docker_args: --net=host {% if not private_dns | is_empty %}--dns={{private_dns}}{% endif %}
    env:
      CONSUL_CLIENT_INTERFACE: "{{consul_bind_interface}}"
      CONSUL_BIND_INTERFACE: "{{consul_bind_interface}}"
```

### Arguments

| Argument       | Default              | Description |
| -------------- | -------------------- | ----------- |
| **image**        | [Required]           | Docker image to run  |
| service | {base image name} | The name of the service (e.g systemd unit name or deployment name) |
| env     |                      | A map of environment variables to pass through |
| labels | | A map of labels to add to the container |
| docker_args |                      | Additional arguments to the docker client e.g. `-p 8080:8080`  |
| docker_opts | | Additional options to the docker client e.g. `-H unix:///tmp/var/run/docker.sock`  |
| args |                   | Additional arguments to the container
| volumes |                | List of volume mappings |
| ports | | List of port mappings |
| commands | | List of commands to execute inside the container on startup ( |
| files | | Map of files to mount into the containers |
| templates | | Map of templates to mount into the container  |
| network | user-bridge |  |
