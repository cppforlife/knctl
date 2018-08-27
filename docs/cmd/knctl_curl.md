## knctl curl

Curl service

### Synopsis

Send a HTTP request to the first ingress address with the Host header set to the service's domain.

Requires 'curl' command installed on the system.

```
knctl curl [flags]
```

### Examples

```

  # Curl service 'svc1' in namespace 'ns1'
  knctl curl -s svc1 -n ns1
```

### Options

```
  -h, --help               help for curl
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -p, --port int32         Set port (default 80)
  -s, --service string     Specified service
```

### Options inherited from parent commands

```
      --column strings      Filter to show only given columns
      --json                Output as JSON
      --kubeconfig string   Path to the kubeconfig file ($KNCTL_KUBECONFIG or $KUBECONFIG)
      --no-color            Disable colorized output
      --non-interactive     Don't ask for user input
      --tty                 Force TTY-like output
```

### SEE ALSO

* [knctl](knctl.md)	 - knctl controls Knative resources (basic-auth-secret, build, curl, deploy, domain, ingress, install, logs, namespace, pod, revision, route, service, service-account, ssh-auth-secret, uninstall, version)

