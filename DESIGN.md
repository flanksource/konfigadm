# Design Principles

## Cloud Native

`konfigadm` has native support for containerized environments and *nothing* else - It is designed to bootstrap immutable container hosts and then handoff to the container(s) and/or orchestrators for everything else.

## Dependency Free

**Runtime Dependencies**
> `konfigadm` is built in pure Go and distributed as a statically linked binary.

Ansible is a good example of how bad dynamically linked tools can get - The dependencies for core Ansible is relatively easy to solve using a Python virtualenv, the 1000's of modules that makeup the extended ansible platform depend on hundreds of other packages - These dependencies are at best visible in documentation and error messages, there is no way of knowing what dependencies (nevermind the version) that are required for a given playbook.

**Explicit Execution Dependencies**
> `konfigadm` does not support any control flows (besides for runtime tags), the ordering of actions is well defined and cannot be changed.

Ansible provides explicit ordering using control statements such as loops and conditionals,  this explicit ordering is coupled with implicit variable management with a very complex precedence hierarchy and ruleset. Creating a mental model of what a given variable will be is almost impossible. Unit testing this state is impossible.

**Implicit Execution Dependencies**
Unlike Ansible, CFEngine and Terraform do not have explicit ordering, Some might argue that they don't have ordering altogether - However the use of input variables and classes create implicit ordering with a complex runtime state machine - It is almost impossible to create and execute this state machine mentally making it very difficult to troubleshoot and test.

## Stateless
> `konfigadm` uses a virtual filesystem and command set for all higher order functions - this makes it trivial to compose and unit-test features.

While the execution model of Ansible does not have persistent state, it is state driven. Facts discovered at runtime can alter behavior, when used across a cluster of machine the impact of state becomes even more important with intermittent connection issues to cluster members potentially creating conflicting state between runs.
The use of change tracking also makes it impossible to ascertain whether a given step will execute.

Ansible, CFEngine and Terraform all execute on the underlying components directly, making only integration testing possible (and difficult)
