#  Templating

## Jinja

`konfigadm` can execute Jinja style templates using the pongo2 library

`config.yml`
```yaml
templates:
  /etc/package-list: file.tpl
```

`file.tpl`
```jinja
{% for pkg in packages %}
Installed package: {{pkg}}
{% endfor %}
```

## YAML tags

### Env

The env tag will fill the environment variable value of `$AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.

```yaml
aws:
  accessKey: !!env AWS_ACCESS_KEY_ID
  secretKey: !!env AWS_SECRET_ACCESS_KEY
```

### Templates

You can use any template function defined by [gomplate](https://github.com/hairyhenderson/gomplate).

```yaml
foo: !!template "{{ base64.Encode \"bar\" }}"
```
