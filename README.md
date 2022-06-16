# ab-example

Sample client-server application that demonstrates use of the Iter8 SDK by the frontend service.

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
