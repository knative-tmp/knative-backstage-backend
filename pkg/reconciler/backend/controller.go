package backend

import (
	"context"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	"k8s.io/client-go/tools/cache"

	eventtypeinformer "knative.dev/eventing/pkg/client/injection/informers/eventing/v1beta1/eventtype"
)

func NewController(ctx context.Context, watcher configmap.Watcher) *controller.Impl {

	reconciler := &Reconciler{
		EventTypeLister: eventtypeinformer.Get(ctx).Lister(),
	}

	logger := logging.FromContext(ctx)

	// TODO:
	logger.Infow("Starting backstage-backend controller")

	impl := brokerreconciler.NewImpl(ctx, reconciler, kafka.BrokerClass, func(impl *controller.Impl) controller.Options {
		return controller.Options{PromoteFilterFunc: kafka.BrokerClassFilter()}
	})

	eventTypeInformer := eventtypeinformer.Get(ctx)

	eventTypeInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			//TODO
			return true
		},
		Handler: controller.HandleAll(impl.Enqueue),
	})

	return impl
}
