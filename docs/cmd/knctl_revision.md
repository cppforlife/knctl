## knctl revision

Revision management (annotate, delete, list, show, tag, untag)

### Synopsis

Revision management (annotate, delete, list, show, tag, untag)

```
knctl revision [flags]
```

### Options

```
  -h, --help   help for revision
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
* [knctl revision annotate](knctl_revision_annotate.md)	 - Annotate revision
* [knctl revision delete](knctl_revision_delete.md)	 - Delete revision
* [knctl revision list](knctl_revision_list.md)	 - List revisions
* [knctl revision show](knctl_revision_show.md)	 - Show revision
* [knctl revision tag](knctl_revision_tag.md)	 - Tag revision
* [knctl revision untag](knctl_revision_untag.md)	 - Untag revision

