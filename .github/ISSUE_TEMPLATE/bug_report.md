---
name: Bug report
about: Create a report to help us improve Capsule
title: ''
labels: blocked-needs-validation, bug
assignees: ''
---

<!--
Thanks for taking time reporting a Nauticus bug!
  
-->

# Bug description

A clear and concise description of what the bug is.

# How to reproduce

Steps to reproduce the behavior:

1. Provide the Nauticus Space YAML definitions
2. Provide all managed Kubernetes resources

# Expected behavior

A clear and concise description of what you expected to happen.

# Logs

If applicable, please provide logs of `nauticus-controller`.

In a standard stand-alone installation of Capsule,
you'd get this by running `kubectl -n nauticus-system logs deploy/nauticus-controller-manager`.

# Additional context

- Nauticus version: 
- Helm Chart version: (`helm list -n nauticus-system`)
- Kubernetes version: (`kubectl version`)