package backend

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"
	eventtypeinformer "knative.dev/eventing/pkg/client/injection/informers/eventing/v1beta2/eventtype"
	eventtypereconciler "knative.dev/eventing/pkg/client/injection/reconciler/eventing/v1beta2/eventtype"
	eventinglistersv1beta2 "knative.dev/eventing/pkg/client/listers/eventing/v1beta2"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"log"
	"net/http"
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

	// TODO: why multiple instances are created?
	handleRoot := func(w http.ResponseWriter, r *http.Request) {
		// TODO: better logging
		fmt.Println("Handling request")

		ret, err := lister.List(labels.Everything())
		if err != nil {
			// TODO: better logs
			log.Fatal(err)
		}

		// TODO
		fmt.Println(ret)

		// write to response as json
		// TODO: hardcoded crap
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(ret)
		if err != nil {
			// TODO: better error handling
			log.Fatal(err)
		}

	}
	http.HandleFunc("/", handleRoot)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
