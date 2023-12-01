package backend

import v1 "knative.dev/eventing/pkg/apis/eventing/v1"

type Trigger struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	UID         string            `json:"uid"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	//
	BrokerRef string `json:"broker,omitempty"`
	//
	ConsumedEventTypes []string `json:"consumedEventTypes,omitempty"`
}

func (tr Trigger) GetNameAndNamespace() string {
	return NameAndNamespace(tr.Namespace, tr.Name)
}

func convertTrigger(tr *v1.Trigger) Trigger {
	var brokerRef = ""
	if tr.Spec.Broker != "" {
		brokerRef = NameAndNamespace(tr.Namespace, tr.Spec.Broker)
	}

	// TODO: use spec.filter
	// TODO: use spec.filters
	// TODO: use spec.destination

	return Trigger{
		Name:        tr.Name,
		Namespace:   tr.Namespace,
		UID:         string(tr.UID),
		Labels:      tr.Labels,
		Annotations: filterAnnotations(tr.Annotations),
		//
		BrokerRef: brokerRef,
		// TODO: check filter
		// ConsumedEventTypes: TODO,
	}
}
