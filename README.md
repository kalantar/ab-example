# ab-example

Sample client-server application that demonstrates use of the Iter8 SDK by the frontend service.

The frontend service provides an endpoint _/hello_ which relies on a backend service endpoint _/world_:

![application interaction](images/application-interaction.png)

To compare multiple versions of the backend service, use the Iter8 SDK to identify which version or _track_ of the backend to send a request to and to export metrics to a metrics data:

![application interation with Iter8 ABn service](images/interaction.png)

An Iter8 experiment can then be written to evaluate the versions.

Sample implementations of the frontend service in go, ... demonstrate the use of the Iter8 API. In these samples, all errors are reported as failures. In practice, a default track might be used.

## Go

```shell
cd go
```

### Build

Set `FRONTEND_TAG` and `BACKEND_TAG` to names of target docker images and build:

```shell
docker build . -f Dockerfile.frontend -t $FRONTEND_TAG
docker push $FRONTEND_TAG
docker build . -f Dockerfile.backend -t $BACKEND_TAG
docker push $BACKEND_TAG
```

### Deploy

Deploy the application:

```shell
sed -e "s#FRONTEND_TAG#$FRONTEND_TAG#" -e "s#BACKEND_TAG#$BACKEND_TAG#" deploy.yaml | kubectl apply -f -
```

### Test

Forward the frontend service:

```shell
kubectl port-forward deploy/frontend 8091:8091
```

Call the application. For example:

```shell
curl localhost:8091/hello -H 'X-User: foo'
```
