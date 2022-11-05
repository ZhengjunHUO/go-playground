package bar

import (
	"context"
	"github.com/ZhengjunHUO/go-playground/pattern/provider-trick/internal"
)

func Setup(mgr *internal.Manager, ctx context.Context) error {
	mgr.Attributes = append(mgr.Attributes, "bar")
	return nil
}
