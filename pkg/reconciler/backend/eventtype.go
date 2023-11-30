package backend

import (
	"context"
	"knative.dev/pkg/logging"
	"net/http"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"

	eventinglistersv1beta2 "knative.dev/eventing/pkg/client/listers/eventing/v1beta2"
)

func EventTypeListHandler(ctx context.Context, lister eventinglistersv1beta2.EventTypeLister) func(w http.ResponseWriter, req *http.Request) {
	logger := logging.FromContext(ctx)

	return func(w http.ResponseWriter, req *http.Request) {
		logger.Debugw("Handling request", "method", req.Method, "url", req.URL)

		ret, err := lister.List(labels.Everything())
		if err != nil {
			logger.Errorw("Error listing eventtypes", "error", err)
			return
		}

		err = json.NewEncoder(w).Encode(ret)
		if err != nil {
			logger.Errorw("Error encoding eventtypes", "error", err)
			return
		}
	}
}
