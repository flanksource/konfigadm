# Integration

`konfigadm` can be used to wrap any other configuration management tool, it  has the advantage of being able to install the configuration management tool, copy the resources required for the tool and execute the tool all in one convenient package.


# Ansible

## Playbooks

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

## Inventory

Ansible inventory files can be used to run `konfigadm` against remote hosts
over SSH. This is supported for the [INI inventory format](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#inventory-basics-formats-hosts-and-groups).
All hosts in the file will be configured and all variables included will be
ignored. Further, hostname ranges, such as `ww[0-5].example.com` are not
supported. `konfigadm` assumes several things:

1. The SSH user is root
1. SSH key authentication uses an SSH agent
1. The SSH port is 22

The inventory can be used by invoking:

```
konfigadm apply -i ./inventory -c config.yml
```
