package backend

import (
	"context"

	v1beta2 "knative.dev/eventing/pkg/apis/eventing/v1beta2"
	eventtypereconciler "knative.dev/eventing/pkg/client/injection/reconciler/eventing/v1beta2/eventtype"
	eventinglistersv1beta2 "knative.dev/eventing/pkg/client/listers/eventing/v1beta2"

	pkgreconciler "knative.dev/pkg/reconciler"
)

type Reconciler struct {
	EventTypeLister eventinglistersv1beta2.EventTypeLister
}

// Check that our Reconciler implements interface
var _ eventtypereconciler.Interface = (*Reconciler)(nil)

func (r *Reconciler) ReconcileKind(ctx context.Context, et *v1beta2.EventType) pkgreconciler.Event {
	// TODO: do we actually need the reconciler?
	return nil
}
