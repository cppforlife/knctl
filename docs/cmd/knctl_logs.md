## knctl logs

Print service logs

### Synopsis

Print service logs of all active pods for a service

```
knctl logs [flags]
```

### Examples

```

  # Fetch last 10 log lines for service 'svc1' in namespace 'ns1' 
  knctl logs -s svc1 -n ns1

  # Follow logs for service 'svc1' in namespace 'ns1' 
  knctl logs -f -s svc1 -n ns1
```

### Options

```
  -f, --follow             As new revisions are added, new pod logs will be printed
  -h, --help               help for logs
  -l, --lines int          Number of lines (default 10)
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -s, --service string     Specified service
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

