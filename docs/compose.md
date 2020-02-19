### Merge Behavior / Composability

Specs can be combined and merged together - e.g. a cloud provider may install PV drivers and a cluster operator may install organization specific motd/issue files.

1. Configuration from files specified later in the chain overwrite previous configurations. (Similar to the ansible [variable precedence](https://docs.ansible.com/ansible/latest/user_guide/playbooks_variables.html#variable-precedence-where-should-i-put-a-variable) rules)
1. **Lists**  are appended to the end of the existing lists (Unsupported in ansible)
1. **Maps** are merged with existing maps (e.g. [hash_behaviour = merge](https://docs.ansible.com/ansible/2.4/intro_configuration.html#hash-behaviour) in ansible)
