package main

import (
	"fmt"
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes/scheme"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		fmt.Printf("[INFO] Found a %s here!\n", val)
		if err := pr.CreateRelatedSvc(ctx, pod); err != nil {
			return reconcile.Result{}, nil
		}
	}

	return reconcile.Result{}, nil
}

func (pr *podReconcile) CreateRelatedSvc(ctx context.Context, pod *corev1.Pod) error {
	if !(len(pod.Spec.Containers) > 0 && len(pod.Spec.Containers[0].Ports) > 0) {
		return nil
	}

	svc := &corev1.Service{}
	svcName := "svc-" + pod.Name
	if err := pr.cl.Get(ctx, client.ObjectKey{
		Namespace: pod.Namespace,
		Name:      svcName,
	}, svc); err != nil {
		if errors.IsNotFound(err) {
			svc = &corev1.Service{
				ObjectMeta: metav1.ObjectMeta {
					Name:      svcName,
					Namespace: pod.Namespace,
				},
				Spec: corev1.ServiceSpec{
					Selector: pod.ObjectMeta.Labels,
					Ports: []corev1.ServicePort {
						{
							Port: pod.Spec.Containers[0].Ports[0].ContainerPort,
						},
					},
					Type: corev1.ServiceTypeNodePort,
				},
			}

			// (owner, controlled metav1.Object, scheme *runtime.Scheme)
			// sets owner as a Controller OwnerReference on controlled
			// used for garbage collection of the controlled object
			// and for reconciling the owner object on changes to controlled (with a Watch + EnqueueRequestForOwner)
			if err = controllerutil.SetControllerReference(pod, svc, scheme.Scheme); err != nil {
				return err
			}

			fmt.Printf("Create service [%s] for pod [%s]\n", svcName, pod.Name)
			return pr.cl.Create(ctx, svc)
		}
		return err
	}

	return nil
}
