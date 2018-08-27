## Annotations

Annotations can be used to store helpful metadata about a service or revision.

```bash
$ export KNCTL_NAMESPACE=default

$ knctl service annotate -a owner=cppforlife --service hello

$ knctl revision annotate -a version=1.3.3 --revision hello:latest

$ knctl revision list --service hello
```
