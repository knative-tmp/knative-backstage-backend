package backend

import (
	eventinglistersv1beta1 "knative.dev/eventing/pkg/client/listers/eventing/v1beta1"
)

type Reconciler struct {
	EventTypeLister eventinglistersv1beta1.EventTypeLister
}
