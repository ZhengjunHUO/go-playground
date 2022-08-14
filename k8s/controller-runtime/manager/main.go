package main

import (
	"fmt"
	"os"
	//"context"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	//"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	corev1 "k8s.io/api/core/v1"
	netwv1 "k8s.io/api/networking/v1"
)

func main() {
	// Init manager
	mngr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		fmt.Println("Create manager failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Manager created.")

	// Init controller, attached to the manager
	ctlrPod, err := controller.New("pod-controller", mngr, controller.Options{
		// Method 1
		//Reconciler: reconcile.Func(func(context.Context, reconcile.Request) (reconcile.Result, error) {
		// 	// Implement reconcile logical here
		//	return reconcile.Result{}, nil
		//}),
		// Method 2
		Reconciler: &podReconcile{cl: mngr.GetClient()},
	})
	if err != nil {
		fmt.Println("Create pod controller failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Pod controller created.")

	// Prepare 2nd controller for svc, attached to the manager
	ctlrSvc, err := controller.New("svc-controller", mngr, controller.Options{
		Reconciler: &svcReconcile{cl: mngr.GetClient()},
	})
	if err != nil {
		fmt.Println("Create svc controller failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Service controller created.")

	// Register handler to controller concerning Pod, add the pod's key to the work queue
	if err := ctlrPod.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}, anotPred); err != nil {
		fmt.Println("Watch pods failed: ", err)
		os.Exit(1)
	}

	// Register handler to controller watching the Service controlled by the pod, add the associated pod's key to the work queue (reconciler only interested in pod in this case)
	// need to call controllerutil.SetControllerReference in the reconcile logic when dealing with the svc related to the pod
	if err := ctlrPod.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &corev1.Pod{},
		IsController: true,
	}); err != nil {
		fmt.Println("Watch svc failed: ", err)
		os.Exit(1)
	}

	// Register handler to controller concerning Service, add the svc's key to the work queue
	if err := ctlrSvc.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForObject{}, anotPred); err != nil {
		fmt.Println("Watch svc failed: ", err)
		os.Exit(1)
	}

	if err := ctlrSvc.Watch(&source.Kind{Type: &netwv1.Ingress{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &corev1.Service{},
		IsController: true,
	}); err != nil {
		fmt.Println("Watch ingress failed: ", err)
		os.Exit(1)
	}

	fmt.Println("Handler registered, run Manager ...")

	// Start manager 
	if err := mngr.Start(signals.SetupSignalHandler()); err != nil {
		fmt.Println("Start manager failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Manager stopped, quit.")
}
