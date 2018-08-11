## knctl

knctl controls Knative resources (annotate, build, create, curl, delete, deploy, install, list, logs, open, route, tag, untag, version)

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

* [knctl annotate](knctl_annotate.md)	 - Annotate (revision, service)
* [knctl build](knctl_build.md)	 - Build source code into image
* [knctl create](knctl_create.md)	 - Create (basic-auth-secret, domain, namespace, service-account, ssh-auth-secret)
* [knctl curl](knctl_curl.md)	 - Curl service
* [knctl delete](knctl_delete.md)	 - Delete (build, revision, route, service)
* [knctl deploy](knctl_deploy.md)	 - Deploy service
* [knctl install](knctl_install.md)	 - Install Knative and Istio
* [knctl list](knctl_list.md)	 - List (builds, domains, ingresses, pods, revisions, routes, services)
* [knctl logs](knctl_logs.md)	 - Print logs
* [knctl open](knctl_open.md)	 - Open web browser pointing at a service domain
* [knctl route](knctl_route.md)	 - Configure route
* [knctl tag](knctl_tag.md)	 - Tag (revision)
* [knctl untag](knctl_untag.md)	 - Untag (revision)
* [knctl version](knctl_version.md)	 - Print client version

