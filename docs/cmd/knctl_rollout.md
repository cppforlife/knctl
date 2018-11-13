## knctl rollout

Create or update route

### Synopsis

Create or update route with traffic percentages.

If route was automatically created for a service, service must be deployed with '--managed-route=false' flag on all subsequent deploys.

```
knctl rollout [flags]
```

### Examples

```

  # Set traffic percentages for service 'svc1' in namespace 'ns1'
  knctl rollout --route rt1 -p svc1:latest=20% -p svc1:previous=80% -n ns1

  # Roll back traffic for previous revision of service 'svc1' in namespace 'ns1'
  knctl rollout --route rt1 -p svc1:previous=100% -n ns1
```

### Options

```
      --generate-name                Set to generate name
  -h, --help                         help for rollout
  -n, --namespace string             Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -p, --percentage strings           Set revision percentage (format: revision=percentage, example: app-00001=100%, app:latest=100%) (can be specified multiple times)
      --route string                 Specified route
      --service-percentage strings   Set service percentage (format: service=percentage, example: app=100%) (can be specified multiple times)
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

* [knctl](knctl.md)	 - knctl controls Knative resources (basic-auth-secret, build, curl, deploy, dns-map, domain, ingress, install, logs, pod, revision, rollout, route, service, service-account, ssh-auth-secret, uninstall, version)

