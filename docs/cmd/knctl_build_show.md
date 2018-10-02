## knctl build show

Show build

### Synopsis

Show build details in a namespace

```
knctl build show [flags]
```

### Examples

```

  # Show details for build 'build1' in namespace 'ns1'
  knctl build show -b build1 -n ns1
```

### Options

```
  -b, --build string       Specified build
  -h, --help               help for show
      --logs               Show logs (default true)
  -n, --namespace string   Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
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

* [knctl build](knctl_build.md)	 - Build management (create, delete, list, show)

