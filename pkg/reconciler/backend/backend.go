package backend

import (
	"context"

	v1beta2 "knative.dev/eventing/pkg/apis/eventing/v1beta2"
	pkgreconciler "knative.dev/pkg/reconciler"
)

type Reconciler struct {
}

func (r *Reconciler) ReconcileKind(ctx context.Context, et *v1beta2.EventType) pkgreconciler.Event {
	// TODO: do we actually need the reconciler?
	return nil
}
