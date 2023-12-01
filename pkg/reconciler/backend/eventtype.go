package backend

import (
	"context"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"
	"knative.dev/eventing/pkg/apis/eventing/v1beta2"
	eventinglistersv1beta2 "knative.dev/eventing/pkg/client/listers/eventing/v1beta2"
	"knative.dev/pkg/logging"
	"net/http"
)

var ExcludedAnnotations = map[string]struct{}{
	"kubectl.kubernetes.io/last-applied-configuration": {},
}

type EventTypeProvider struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
}

type EventType struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Type        string            `json:"type"`
	UID         string            `json:"uid"`
	Description string            `json:"description,omitempty"`
	SchemaData  string            `json:"schemaData,omitempty"`
	SchemaURL   string            `json:"schemaURL,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Provider    EventTypeProvider `json:"provider,omitempty"`
}

func EventTypeListHandler(ctx context.Context, lister eventinglistersv1beta2.EventTypeLister) func(w http.ResponseWriter, req *http.Request) {
	logger := logging.FromContext(ctx)

	return func(w http.ResponseWriter, req *http.Request) {
		logger.Debugw("Handling request", "method", req.Method, "url", req.URL)

		ret, err := lister.List(labels.Everything())
		if err != nil {
			logger.Errorw("Error listing eventtypes", "error", err)
			return
		}

		converted := make([]EventType, 0, len(ret))
		for _, et := range ret {
			converted = append(converted, convertEventType(et))
		}

		err = json.NewEncoder(w).Encode(converted)
		if err != nil {
			logger.Errorw("Error encoding eventtypes", "error", err)
			return
		}
	}
}

func convertEventType(et *v1beta2.EventType) EventType {
	provider := EventTypeProvider{}
	if et.Spec.Reference != nil {
		provider.Name = et.Spec.Reference.Name
		provider.Kind = et.Spec.Reference.Kind
		provider.Namespace = et.Spec.Reference.Namespace
	}

	// TODO: need deduplication at this level here

	// TODO: more information!
	return EventType{
		Name:        et.Name,
		Namespace:   et.Namespace,
		Type:        et.Spec.Type,
		UID:         string(et.UID),
		Description: et.Spec.Description,
		SchemaData:  et.Spec.SchemaData,
		SchemaURL:   et.Spec.Schema.String(),
		Labels:      et.Labels,
		Annotations: filterAnnotations(et.Annotations),
		// TODO: remove?
		Provider: provider,
	}
}

func filterAnnotations(annotations map[string]string) map[string]string {
	ret := make(map[string]string)
	for k, v := range annotations {
		if _, ok := ExcludedAnnotations[k]; ok {
			continue
		}
		ret[k] = v
	}
	return ret
}
