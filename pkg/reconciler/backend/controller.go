package backend

import (
	"context"
	eventtypereconciler "knative.dev/eventing/pkg/client/injection/reconciler/eventing/v1beta2/eventtype"

	brokerinformer "knative.dev/eventing/pkg/client/injection/informers/eventing/v1/broker"
	triggerinformer "knative.dev/eventing/pkg/client/injection/informers/eventing/v1/trigger"
	eventtypeinformer "knative.dev/eventing/pkg/client/injection/informers/eventing/v1beta2/eventtype"

	eventinglistersv1 "knative.dev/eventing/pkg/client/listers/eventing/v1"
	eventinglistersv1beta2 "knative.dev/eventing/pkg/client/listers/eventing/v1beta2"

	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Listers struct {
	EventTypeLister eventinglistersv1beta2.EventTypeLister
	BrokerLister    eventinglistersv1.BrokerLister
	TriggerLister   eventinglistersv1.TriggerLister
}

func NewController(ctx context.Context) *controller.Impl {

	reconciler := &Reconciler{}

	logger := logging.FromContext(ctx)

	logger.Infow("Starting backstage-backend controller")

	impl := eventtypereconciler.NewImpl(ctx, reconciler)

	listers := Listers{
		EventTypeLister: eventtypeinformer.Get(ctx).Lister(),
		BrokerLister:    brokerinformer.Get(ctx).Lister(),
		TriggerLister:   triggerinformer.Get(ctx).Lister(),
	}
	go startWebServer(ctx, listers)

	return impl
}

func startWebServer(ctx context.Context, listers Listers) {

	logger := logging.FromContext(ctx)

	logger.Infow("Starting backstage-backend webserver")

	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.HandleFunc("/", EventMeshHandler(ctx, listers)).Methods("GET")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
