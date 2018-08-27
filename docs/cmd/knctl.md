## knctl

knctl controls Knative resources (basic-auth-secret, build, curl, deploy, domain, ingress, install, logs, namespace, pod, revision, route, service, service-account, ssh-auth-secret, uninstall, version)

### Synopsis

knctl controls Knative resources.

CLI docs: https://github.com/cppforlife/knctl#docs.
Knative docs: https://github.com/knative/docs.

```
knctl [flags]
```

### Options

```
      --column strings      Filter to show only given columns
  -h, --help                help for knctl
      --json                Output as JSON
      --kubeconfig string   Path to the kubeconfig file ($KNCTL_KUBECONFIG) (default "/Users/pivotal/.kube/config")
      --no-color            Disable colorized output
      --non-interactive     Don't ask for user input
      --tty                 Force TTY-like output
```

### SEE ALSO

* [knctl basic-auth-secret](knctl_basic-auth-secret.md)	 - Basic auth secret (create)
* [knctl build](knctl_build.md)	 - Build (create, delete, list)
* [knctl curl](knctl_curl.md)	 - Curl service
* [knctl deploy](knctl_deploy.md)	 - Deploy service
* [knctl domain](knctl_domain.md)	 - Domain (create, list)
* [knctl ingress](knctl_ingress.md)	 - Ingress (list)
* [knctl install](knctl_install.md)	 - Install Knative and Istio
* [knctl logs](knctl_logs.md)	 - Print logs
* [knctl namespace](knctl_namespace.md)	 - Namespace (create)
* [knctl pod](knctl_pod.md)	 - Pod (list)
* [knctl revision](knctl_revision.md)	 - Revision (annotate, delete, list, tag, untag)
* [knctl route](knctl_route.md)	 - Route (create, delete, list)
* [knctl service](knctl_service.md)	 - Service (annotate, delete, list, open)
* [knctl service-account](knctl_service-account.md)	 - Service account (create)
* [knctl ssh-auth-secret](knctl_ssh-auth-secret.md)	 - SSH auth secret (create)
* [knctl uninstall](knctl_uninstall.md)	 - Uninstall Knative and Istio
* [knctl version](knctl_version.md)	 - Print client version

