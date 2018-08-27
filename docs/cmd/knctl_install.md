## knctl install

Install Knative and Istio

### Synopsis

Install Knative and Istio.

Requires 'kubectl' command installed on a the system.

```
knctl install [flags]
```

### Options

```
  -m, --exclude-monitoring   Exclude installation of monitoring components
  -h, --help                 help for install
  -p, --node-ports           Use service type NodePorts instead of type LoadBalancer
      --version-check        Check minimum Kubernetes API server version (default true)
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

