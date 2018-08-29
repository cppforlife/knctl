## knctl build delete

Delete build

### Synopsis

Delete build

```
knctl build delete [flags]
```

### Examples

```

  # Delete build 'build1' in namespace 'ns1'
  knctl build delete -b build1 -n ns1
```

### Options

```
  -b, --build string       Specified build
  -h, --help               help for delete
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

* [knctl build](knctl_build.md)	 - Build management (create, delete, list)

