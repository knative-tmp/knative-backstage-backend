package main

import (
	"context"
	"knative.dev/backstage-backend/pkg/reconciler/backend"
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

		injection.NamedControllerConstructor{
			Name: "backend",
			ControllerConstructor: func(ctx context.Context, watcher configmap.Watcher) *controller.Impl {
				return backend.NewController(ctx)
			},
		},
	)
}
