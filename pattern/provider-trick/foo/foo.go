package foo

import (
	"context"
	"github.com/ZhengjunHUO/playground/pattern/provider-trick/internal"
)

func Setup(mgr *internal.Manager, ctx context.Context) error {
	mgr.Attributes = append(mgr.Attributes, "foo")
	return nil
}
