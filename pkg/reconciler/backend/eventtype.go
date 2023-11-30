package backend

import (
	"fmt"
	"log"
	"net/http"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"

	eventinglistersv1beta2 "knative.dev/eventing/pkg/client/listers/eventing/v1beta2"
)

func EventTypeListHandler(lister eventinglistersv1beta2.EventTypeLister) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
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
}
