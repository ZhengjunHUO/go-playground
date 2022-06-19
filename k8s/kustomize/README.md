## Deployment
```bash
# Create apps in prod namespace
kubectl apply -k ./prod/
# Create apps in dev namespace
kubectl apply -k ./dev/
```
## Check
```bash
kubectl run alpine --image=alpine --restart=Never --command -- sleep infinity
kubectl exec -ti alpine -- sh
/ # apk add curl
/ # curl http-server-ep.prod.svc:8080
Welcome to prod env.
/ # curl http-server-ep.dev.svc:8080
Welcome to dev env.
```

## Misc
```bash
# if deploy the base directory, will create the apps in the current namespace
kubectl apply -k ./base/
```
