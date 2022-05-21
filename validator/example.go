package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type k8sCluster struct {
	Name		string	`validate:"required,alphanum,max=20,min=3"`
	CNI		string	`validate:"required,oneof=cilium calico flannel weave"`
	Size		int	`validate:"required,numeric,min=1"`
	IsManaged	bool	`validate:"omitempty"`
	IsBaremetal	bool	`validate:"omitempty"`
	IsOverlay	bool	`validate:"omitempty"`
}

func main() {
	vldt := validator.New()

	cluster := k8sCluster{
		Name:		"huo",
		Size:		8,
		CNI:		"cilium",
		IsManaged:	false,
		IsBaremetal:	true,
		IsOverlay:	false,
	}

	if err := vldt.Struct(cluster); err != nil {
		fmt.Println(err)
	}else{
		fmt.Println("Validation ok.")
	}
}
