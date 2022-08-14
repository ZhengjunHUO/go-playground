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
	netwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type podReconcile struct {
	cl client.Client
}

type svcReconcile struct {
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
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (sr *svcReconcile) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	fmt.Printf("Svc: [%s] in ns [%s]\n", req.Name, req.Namespace)

	svc := &corev1.Service{}
	if err := sr.cl.Get(ctx, req.NamespacedName, svc); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, fmt.Errorf("Get svc failed: %v\n", err)
	}

	if val, ok := svc.Annotations["huozj.io/animals"]; ok && val == "cat" {
		fmt.Printf("[INFO] Found a %s in svc [%s/%s]!\n", val, svc.Namespace, svc.Name)
		if err := sr.CreateRelatedIngress(ctx, svc); err != nil {
			return reconcile.Result{}, err
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
					Name:        svcName,
					Namespace:   pod.Namespace,
					Annotations: pod.ObjectMeta.Annotations,
				},
				Spec: corev1.ServiceSpec{
					Selector: pod.ObjectMeta.Labels,
					Ports: []corev1.ServicePort {
						{
							Port: pod.Spec.Containers[0].Ports[0].ContainerPort,
						},
					},
					Type: corev1.ServiceTypeClusterIP,
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

func (sr *svcReconcile) CreateRelatedIngress(ctx context.Context, svc *corev1.Service) error {
	if len(svc.Spec.Ports) == 0 {
		return nil
	}

	ing := &netwv1.Ingress{}
	ingName := "ing-" + svc.Name
	if err := sr.cl.Get(ctx, client.ObjectKey{
		Namespace: svc.Namespace,
		Name:      ingName,
	}, ing); err != nil {
		if errors.IsNotFound(err) {
			prefix := netwv1.PathTypePrefix
			ing = &netwv1.Ingress{
				ObjectMeta: metav1.ObjectMeta {
					Name:        ingName,
					Namespace:   svc.Namespace,
					Annotations: map[string]string{
						"kubernetes.io/ingress.class": "nginx",
					},
				},
				Spec: netwv1.IngressSpec{
					Rules: []netwv1.IngressRule{
						{
							Host: "www.huozj.io",
							IngressRuleValue: netwv1.IngressRuleValue{
								HTTP: &netwv1.HTTPIngressRuleValue{
									Paths: []netwv1.HTTPIngressPath{
										{
											Path: "/",
											PathType: &prefix,
											Backend: netwv1.IngressBackend{
												Service: &netwv1.IngressServiceBackend{
													Name: svc.Name,
													Port: netwv1.ServiceBackendPort{Number: svc.Spec.Ports[0].Port,},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			if err = controllerutil.SetControllerReference(svc, ing, scheme.Scheme); err != nil {
				fmt.Printf("Set controller ref failed for [%s]: %v\n", ingName, err)
				return err
			}

			if err := sr.cl.Create(ctx, ing); err != nil {
				fmt.Printf("Get an error in creating ingress [%s]: %v\n", ingName, err)
				return err
			}

			fmt.Printf("Create ingress [%s] for svc [%s]\n", ingName, svc.Name)
		}
		return err
	}

	return nil
}
