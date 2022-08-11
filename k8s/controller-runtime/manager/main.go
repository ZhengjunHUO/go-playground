package main

import (
	"fmt"
	"os"
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	corev1 "k8s.io/api/core/v1"
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
	ctlr, err := controller.New("demo-controller", mngr, controller.Options{
		Reconciler: reconcile.Func(func(context.Context, reconcile.Request) (reconcile.Result, error) {
			return reconcile.Result{}, nil
		}),
	})
	if err != nil {
		fmt.Println("Create controller failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Controller created.")

	// Register handler to controller concerning Pod, add to the work queue
	if err := ctlr.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{}); err != nil {
		fmt.Println("Watch pods failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Handler registered, run Manager ...")

	// Start manager 
	if err := mngr.Start(signals.SetupSignalHandler()); err != nil {
		fmt.Println("Watch pods failed: ", err)
		os.Exit(1)
	}
	fmt.Println("Manager stopped, quit.")
}
