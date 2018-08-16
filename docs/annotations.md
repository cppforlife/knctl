## Annotations

Annotations can be used to store helpful metadata about a service or revision.

```bash
$ export KNCTL_NAMESPACE=default

$ knctl annotate service -a owner=cppforlife --service hello

$ knctl annotate revision -a version=1.3.3 --revision hello:latest

$ knctl list revisions --service hello
```
