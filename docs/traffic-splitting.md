## Traffic Splitting WIP

**Currently does not work with latest Knative Serving**

### Tags

Each tag identifies single revision within a service. CLI uses Kubernetes labels to store tag information. Tags can be used to reference particular revision when changing traffic configuration. By default, CLI will apply two tags: `latest` and `previous`.

```bash
$ knctl -n default revision tag --revision hello:latest -t stable

$ knctl -n default revision list --service hello
```

### Routing

`knctl` has not yet fully implemented traffic splitting functionality available in Knative Serving. In many situations, developers want to expose only a percentage of total traffic to particular version of their application. This may be useful to be able to verify popularity of some features, to test out aesthetics of a user interface, or to gradually verify that a fix actually resolves an issue. Traffic splitting may help with these use cases.

In a nutshell, developer wants to route different percentage of the traffic to different revisions. For instance, to verify a fix, a developer might want to have 10 to 20% percent traffic routed to the version with fix. Monitor traffic to the end-point where the fix is best observed and see if the fix caused more problems than solved.

After enough data is gathered, developers can then route 100% traffic to the select revision. As you can see the actual percetage and number of revisions invloved changes based on the actual use case in question. For that reason we want to have a `knctl` experience that is flexible for different use cases. The gist of the experience looks like this command:

- to direct a bit more of traffic to new revision

```bash
$ knctl route create --route hello -p hello:latest=40% -p hello:previous=60%
```

- to direct all traffic to latest revision

```bash
$ knctl route create --route hello -p hello:latest=100%
```

- to direct all traffic back to previous revision

```bash
$ knctl route create --route hello -p hello:previous=100%
```
