package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/locales/en"
        ut "github.com/go-playground/universal-translator"
	enTrans "github.com/go-playground/validator/v10/translations/en"
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

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")

	_ = enTrans.RegisterDefaultTranslations(vldt, trans)

	clusters := []k8sCluster{
	    k8sCluster{
		Name:		"x",
		Size:		0,
		CNI:		"awsvps",
		IsManaged:	false,
	    },
	    k8sCluster{
		Name:		"huo",
		Size:		8,
		CNI:		"cilium",
		IsManaged:	false,
		IsBaremetal:	true,
		IsOverlay:	false,
	    },
	}

	for i := range clusters {
		fmt.Printf("<Cluster %d>\n", i)
		if err := vldt.Struct(clusters[i]); err != nil {
			fmt.Println(interpretError(err.(validator.ValidationErrors), trans))
		}else{
			fmt.Println("Validation ok.")
		}
	}
}

func interpretError(err validator.ValidationErrors, trans ut.Translator) []error {
	errs := []error{}
	if err != nil {
		for _, e := range err {
			errs = append(errs, fmt.Errorf("%s\n", e.Translate(trans)))
		}
	}

	return errs
}
