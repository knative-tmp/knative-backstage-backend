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

	go startWebServer(ctx, eventtypeinformer.Get(ctx).Lister())

	return impl
}

func startWebServer(ctx context.Context, lister eventinglistersv1beta2.EventTypeLister) {

	logger := logging.FromContext(ctx)

	logger.Infow("Starting backstage-backend webserver")

	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.HandleFunc("/eventtypes", EventTypeListHandler(ctx, lister)).Methods("GET")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
