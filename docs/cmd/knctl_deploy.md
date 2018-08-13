## knctl deploy

Deploy service

### Synopsis

Deploy service

```
knctl deploy [flags]
```

### Examples

```

  # Deploy service 'srv1' with a given image and one environment variable in namespace 'ns1'
  knctl deploy -s srv1 --image gcr.io/knative-samples/helloworld-go --env TARGET=123 -n ns1

  # Deploy service 'srv1' from Git repo and one environment variable in namespace 'ns1'
  knctl deploy -s srv1 --image gcr.io/your-account/your-image --git-url https://github.com/cppforlife/simple-app --git-revision master --env TARGET=123 -n ns1
```

### Options

```
  -e, --env strings              Set environment variable (format: key=value) (can be specified multiple times)
      --generate-name            Set to generate name
      --git-revision string      Set Git revision (examples: https://git-scm.com/docs/gitrevisions#_specifying_revisions)
      --git-url string           Set Git URL
  -h, --help                     help for deploy
  -i, --image string             Set image URL
  -n, --namespace string         Specified namespace ($KNCTL_NAMESPACE)
  -s, --service string           Specified service
      --service-account string   Set service account name for building
```

### Options inherited from parent commands

```
      --column strings      Filter to show only given columns
      --json                Output as JSON
      --kubeconfig string   Path to the kubeconfig file ($KNCTL_KUBECONFIG) (default "/Users/pivotal/.kube/config")
      --no-color            Disable colorized output
      --non-interactive     Don't ask for user input
      --tty                 Force TTY-like output
```

### SEE ALSO

* [knctl](knctl.md)	 - knctl controls Knative resources (annotate, build, create, curl, delete, deploy, install, list, logs, open, route, tag, untag, version)

