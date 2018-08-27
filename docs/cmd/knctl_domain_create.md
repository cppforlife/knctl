## knctl domain create

Create domain

### Synopsis

Create domain

```
knctl domain create [flags]
```

### Examples

```

  # Create domain 'example.com' and set it as default
  knctl domain create -d example.com --default
```

### Options

```
      --default         Set domain as default (currently required to be provided)
  -d, --domain string   Specified domain (example: domain.com)
  -h, --help            help for create
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

* [knctl domain](knctl_domain.md)	 - Domain (create, list)

