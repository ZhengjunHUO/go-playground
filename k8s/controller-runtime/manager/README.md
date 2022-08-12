```sh
$ kubectl run nginx --image=nginx --port 80
$ kubectl annotate pod nginx huozj.io/animals=cat
# Check the service created by controller
$ kubectl get svc
# Delete the service and redo get to see what happened
# The service is watched (by 2nd handler) and will be recreated by the reconcile logic
$ kubectl delete svc svc-nginx
# Delete the pod and check again the service
# The garbage collector remove also the svc, since the owner pod is deleted 
$ kubectl delete pod nginx
```
