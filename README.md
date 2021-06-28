# config-connector-templater

Kubernetes operator that templates GCP Config Connector resources

## Example

Creating a resource like this:

```yaml
apiVersion: config-connector-templater.slamdev.net/v1alpha1
kind: PubSubTopicTemplate
metadata:
  name: notifications
  namespace: team1
  annotations:
    service-name: super-service
spec:
  resourceID: "{{ .metadata.namespace }}.{{ index .metadata.annotations "service-name" }}.{{ .metadata.name }}"
```

operator will create the following resource:

```yaml
apiVersion: pubsub.cnrm.cloud.google.com/v1beta1
kind: PubSubTopic
metadata:
  name: notifications
spec:
  resourceID: team1.super-service.notifications
```

## Make a release

```shell script
TAG=x.x.x && git tag -a ${TAG} -m "make ${TAG} release" && git push --tags
```
