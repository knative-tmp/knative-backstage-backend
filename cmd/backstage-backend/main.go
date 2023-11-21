package main

import (
	"context"

	"github.com/knative-tmp/knative-backstage-backend/pkg/reconciler/backend"

	"knative.dev/pkg/injection"
	"knative.dev/pkg/signals"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/sharedmain"
)

const (
	component = "backstage-backend"
)

func main() {

	sharedmain.MainNamed(signals.NewContext(), component,

		// Broker controller
		injection.NamedControllerConstructor{
			Name: "backend",
			ControllerConstructor: func(ctx context.Context, watcher configmap.Watcher) *controller.Impl {
				return backend.NewController(ctx, watcher)
			},
		},
	)
}
