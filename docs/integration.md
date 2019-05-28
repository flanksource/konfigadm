# Integration

`konfigadm` can be used to wrap any other configuration management tool, it  has the advantage of being able to install the configuration management tool, copy the resources required for the tool and execute the tool all in one convenient package.


## ansible

`config.yml`
```yaml
packages:
  - ansible
files:
  /root/playbook.yml: |
    ---
    - name: This is a hello-world example
      hosts: all
      connection: local
      tasks:
        - name: Create a file called '/tmp/testfile.txt'
          copy:
            content: hello world
            dest: /tmp/testfile.txt
commands:
  - ansible-playbook -i 'localhost, ' /root/playbook.yml
```
```bash
konfigadm apply -c config.yml
```

The playbook and other files can also be externalized by just specifying a relative path to the files to include.

`config.yml`
```yaml
packages:
  - ansible
files:
  /root/playbook.yml: playbook.yml
commands:
  - ansible-playbook -i 'localhost, ' /root/playbook.yml
```

`playbook.yml`
```yaml
- name: This is a hello-world example
    hosts: all
    connection: local
    tasks:
      - name: Create a file called '/tmp/testfile.txt'
        copy:
          content: hello world
          dest: /tmp/testfile.txt
```
