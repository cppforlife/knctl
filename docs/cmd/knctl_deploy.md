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
  knctl deploy -s srv1 --image gcr.io/your-account/your-image \
      --git-url https://github.com/cppforlife/simple-app --git-revision master --env TARGET=123 -n ns1

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

  # Deploy service 'srv1' that needs secret values in environment variables
  # ( https://github.com/cppforlife/knctl/blob/master/docs/deploy-secrets.md )
  knctl deploy -s srv1 -n ns1 \
      --image gcr.io/knative-samples/helloworld-go \
      --env-secret TARGET=secret/key1 \
      --env-secret TARGET=secret/key2
	  
  knctl deploy -s srv1 -n ns1 \
      --image gcr.io/knative-samples/helloworld-go \
	  --secret-mount secret-name=/mount/path1
	  --config-map-mount config-map-name=/mount/path2
	  
```

### Options

```
  -a, --annotation strings                      Set annotation (format: key=value) (can be specified multiple times)
      --build-arg stringArray                   Set build argument (format: key=value) (can be specified multiple times)
      --build-timeout duration                  Set timeout for building stage (Knative Build has a 10m default)
      --config-map-mount strings                Mount a config map as a volume (format: configmap-name=/mount/path) (can be specified multiple times)
      --container-concurrency int               Set container concurrency (default unspecified)
  -d, --directory string                        Set source code directory
      --dry-run                                 Dry run
  -e, --env stringArray                         Set environment variable (format: ENV_KEY=value) (can be specified multiple times)
      --env-all-from-config-map strings         Set environment variables as all key-value in a config map (format: config-map-name) (can be specified multiple times)
      --env-config-map strings                  Set environment variable from a config map (format: ENV_KEY=config-map-name/key) (can be specified multiple times)
      --env-secret strings                      Set environment variable from a secret (format: ENV_KEY=secret-name/key) (can be specified multiple times)
      --generate-name                           Set to generate name
      --git-revision string                     Set Git revision (examples: https://git-scm.com/docs/gitrevisions#_specifying_revisions)
      --git-url string                          Set Git URL
  -h, --help                                    help for deploy
  -i, --image string                            Set image URL
      --managed-route                           Custom route configuration (default true)
      --max-scale int                           Set autoscaling rule for maximum number of containers (default unspecified)
      --min-scale int                           Set autoscaling rule for minimum number of containers (default unspecified)
  -n, --namespace string                        Specified namespace ($KNCTL_NAMESPACE or default from kubeconfig)
      --secret-mount strings                    Mount a secret as a volume (format: secret-name=/mount/path) (can be specified multiple times)
  -s, --service string                          Specified service
      --service-account string                  Set service account name for building
  -t, --tag strings                             Set tag (format: value) (can be specified multiple times)
      --template string                         Set template name
      --template-arg stringArray                Set template argument (format: key=value) (can be specified multiple times)
      --template-env stringArray                Set template environment variable (format: key=value) (can be specified multiple times)
      --template-kind string                    Set to 'cluster' to use ClusterBuildTemplate kind of templates
      --watch-pod-logs                          Watch pod logs for new revision (default true)
  -l, --watch-pod-logs-indefinitely             Watch pod logs for new revision indefinitely
      --watch-revision-ready                    Wait for new revision to become ready (default true)
      --watch-revision-ready-timeout duration   Set timeout for waiting for new revision to become ready (default 5m0s)
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

