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

  # Deploy service 'srv1' from local source code in namespace 'ns1'
  # ( https://github.com/cppforlife/knctl/blob/master/docs/deploy-source-directory.md )
  knctl deploy -s srv1 -d=. --image index.docker.io/your-account/your-image --service-account serv-acct1 --env TARGET=123 -n ns1

  # Deploy service 'srv1' with custom build template in namespace 'ns1'
  # ( https://github.com/cppforlife/knctl/blob/master/docs/deploy-custom-build-template.md )
  knctl deploy -s srv1 -n ns1 \
      --git-url https://github.com/cppforlife/simple-app --git-revision master \
      --template buildpack --template-env GOPACKAGENAME=main \
      --service-account serv-acct1 --image index.docker.io/your-account/your-repo \
      --env SIMPLE_MSG=123
```

### Options

```
  -c, --cluster-registry                    Use cluster registry
      --cluster-registry-namespace string   Namespace where cluster registry was installed
  -d, --directory string                    Set source code directory
  -e, --env strings                         Set environment variable (format: key=value) (can be specified multiple times)
      --generate-name                       Set to generate name
      --git-revision string                 Set Git revision (examples: https://git-scm.com/docs/gitrevisions#_specifying_revisions)
      --git-url string                      Set Git URL
  -h, --help                                help for deploy
  -i, --image string                        Set image URL
  -n, --namespace string                    Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
  -s, --service string                      Specified service
      --service-account string              Set service account name for building
      --template string                     Set template name
      --template-arg strings                Set template argument (format: key=value) (can be specified multiple times)
      --template-env strings                Set template environment variable (format: key=value) (can be specified multiple times)
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

