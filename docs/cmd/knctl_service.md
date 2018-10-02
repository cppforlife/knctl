## knctl service

Service management (annotate, delete, list, open, show, url)

### Synopsis

Service management (annotate, delete, list, open, show, url)

```
knctl service [flags]
```

### Options

```
  -h, --help   help for service
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

* [knctl](knctl.md)	 - knctl controls Knative resources (basic-auth-secret, build, curl, deploy, dns-map, domain, ingress, install, logs, namespace, pod, revision, route, service, service-account, ssh-auth-secret, uninstall, version)
* [knctl service annotate](knctl_service_annotate.md)	 - Annotate service
* [knctl service delete](knctl_service_delete.md)	 - Delete service
* [knctl service list](knctl_service_list.md)	 - List services
* [knctl service open](knctl_service_open.md)	 - Open web browser pointing at a service domain
* [knctl service show](knctl_service_show.md)	 - Show service
* [knctl service url](knctl_service_url.md)	 - Print service URL

