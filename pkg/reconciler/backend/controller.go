package backend

import (
	"context"
	eventtypeinformer "knative.dev/eventing/pkg/client/injection/informers/eventing/v1beta2/eventtype"
	eventtypereconciler "knative.dev/eventing/pkg/client/injection/reconciler/eventing/v1beta2/eventtype"
	eventinglistersv1beta2 "knative.dev/eventing/pkg/client/listers/eventing/v1beta2"

	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewController(ctx context.Context) *controller.Impl {

	reconciler := &Reconciler{
		EventTypeLister: eventtypeinformer.Get(ctx).Lister(),
	}

	logger := logging.FromContext(ctx)

	logger.Infow("Starting backstage-backend controller")

	impl := eventtypereconciler.NewImpl(ctx, reconciler)

	go startWebServer(eventtypeinformer.Get(ctx).Lister())

	return impl
}

func startWebServer(lister eventinglistersv1beta2.EventTypeLister) {

	r := mux.NewRouter()
	r.HandleFunc("/eventtypes", EventTypeListHandler(lister))
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
