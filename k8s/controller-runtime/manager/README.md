### [Note] annotate the pod will trigger the controller to create a associated service, then another controller who is watching the service will create a assocaite ingress (chained ownerReferences)
```sh
$ kubectl run nginx --image=nginx --port 80
$ kubectl annotate pod nginx huozj.io/animals=cat
# Check the service created by controller and ingress created by service
$ kubectl get svc,ing
# Delete the service/ingress and redo get to see what happened
# The service is watched (by 2nd handler) and will be recreated by the reconcile logic; idem for ingress
$ kubectl delete svc svc-nginx
# Delete the pod and check again the service and ingress
# The garbage collector remove also the svc, since the owner pod is deleted; idem for ingress
$ kubectl delete pod nginx
```
