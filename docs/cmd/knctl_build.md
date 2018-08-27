## knctl build

Build (create, delete, list)

### Synopsis

Build (create, delete, list)

```
knctl build [flags]
```

### Options

```
  -h, --help   help for build
```

### Options inherited from parent commands

```
      --column strings              Filter to show only given columns
      --json                        Output as JSON
      --kubeconfig string           Path to the kubeconfig file ($KNCTL_KUBECONFIG or $KUBECONFIG)
      --kubeconfig-context string   Kubeconfig context override ($KNCTL_KUBECONFIG_CONTEXT)
      --no-color                    Disable colorized output
      --non-interactive             Don't ask for user input
      --tty                         Force TTY-like output
```

### SEE ALSO

* [knctl](knctl.md)	 - knctl controls Knative resources (basic-auth-secret, build, curl, deploy, domain, ingress, install, logs, namespace, pod, revision, route, service, service-account, ssh-auth-secret, uninstall, version)
* [knctl build create](knctl_build_create.md)	 - Build source code into image
* [knctl build delete](knctl_build_delete.md)	 - Delete build
* [knctl build list](knctl_build_list.md)	 - List builds

