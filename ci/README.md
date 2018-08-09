## ci

pipeline.yml is provided for running tests in Concourse:

```
$ fly -t env set-pipeline -p knctl -c ci/pipeline.yml -l config.yml
```

where `config.yml` is in following format:

```
kubeconfig: |
  apiVersion: v1
  clusters:
  - cluster:
    ...
```

Currently, it's expected that Knative is pre-installed into specified cluster.
