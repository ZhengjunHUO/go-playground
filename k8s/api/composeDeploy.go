package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
)

func main() {
	var two int32 = 2
	deploy := appsv1.Deployment {
		Spec: appsv1.DeploymentSpec {
			Replicas: &two,
			Template: corev1.PodTemplateSpec {
				Spec: corev1.PodSpec {
					Containers: []corev1.Container {
						{Name: "test", Image: "alpine", Command: []string{"sleep infinity"}},
					},
				},
			},
		},
	}

	fmt.Printf("Deployment test: %#v\n", &deploy)
}
