package main

import (
	"fmt"
	"context"

	"github.com/ZhengjunHUO/go-playground/pattern/provider-trick/internal"
	"github.com/ZhengjunHUO/go-playground/pattern/provider-trick/foo"
	"github.com/ZhengjunHUO/go-playground/pattern/provider-trick/bar"
)

func main() {
	ctx := context.Background()
	mgr := &internal.Manager{
		Attributes: []string{},
	}

	if err := Setup(mgr, ctx); err != nil {
		fmt.Println(err)
	}else{
		fmt.Println(mgr.Attributes)
	}
}

func Setup(mgr *internal.Manager, ctx context.Context) error {
	for _, setup := range []func(mgr *internal.Manager, ctx context.Context) error {
		foo.Setup,
		bar.Setup,
	} {
		if err := setup(mgr, ctx); err != nil {
			return err
		}
	}

	return nil
}
