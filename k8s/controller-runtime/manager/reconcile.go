package main

import (
	"fmt"
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"
)

type podReconcile struct {
	cl client.Client
}

func (pr *podReconcile) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	fmt.Printf("Pod: [%s] in ns [%s]\n", req.Name, req.Namespace)

	pod := &corev1.Pod{}
	if err := pr.cl.Get(ctx, req.NamespacedName, pod); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, fmt.Errorf("Get pod failed: %v\n", err)
	}

	if val, ok := pod.Annotations["huozj.io/animals"]; ok && val == "cat" {
		// TO IMPLEMENT
		fmt.Printf("Found a %s here!\n", val)
	}

	return reconcile.Result{}, nil
}
