#  Templating

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
