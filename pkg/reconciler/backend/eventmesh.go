package backend

import (
	"context"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/json"
	eventinglistersv1 "knative.dev/eventing/pkg/client/listers/eventing/v1"
	"knative.dev/pkg/logging"
	"net/http"
)

type EventMesh struct {
	// not every event type is tied to a broker. thus, we need to send event types as well.
	EventTypes []EventType `json:"eventTypes"`
	Brokers    []*Broker   `json:"brokers"`
	Triggers   []Trigger   `json:"triggers"`
}

type EventTypeMap = map[string]*EventType

func EventMeshHandler(ctx context.Context, listers Listers) func(w http.ResponseWriter, req *http.Request) {
	logger := logging.FromContext(ctx)

	return func(w http.ResponseWriter, req *http.Request) {
		logger.Debugw("Handling request", "method", req.Method, "url", req.URL)

		err, convertedBrokers := fetchBrokers(listers.BrokerLister, logger)
		if err != nil {
			logger.Errorw("Error fetching and converting brokers", "error", err)
			return
		}

		brokerMap := make(map[string]*Broker)
		for _, cbr := range convertedBrokers {
			brokerMap[cbr.GetNameAndNamespace()] = cbr
		}

		fetchedEventTypes, err := listers.EventTypeLister.List(labels.Everything())
		if err != nil {
			logger.Errorw("Error listing eventTypes", "error", err)
			return
		}

		convertedEventTypes := make([]EventType, 0, len(fetchedEventTypes))
		for _, et := range fetchedEventTypes {
			convertedEventTypes = append(convertedEventTypes, convertEventType(et))

			if et.Spec.Reference != nil {
				if br, ok := brokerMap[RefNameAndNamespace(et.Spec.Reference)]; ok {
					br.ProvidedEventTypes = append(br.ProvidedEventTypes, ObjNameAndNamespace(et))
				}
			}
		}

		//brokerEventTypeMap := make(map[string]EventTypeMap)
		//
		//for _, et := range convertedEventTypes {
		//	if _, ok := brokerEventTypeMap[et.]; !ok {
		//		brokerEventTypeMap[et.BrokerRef] = make(EventTypeMap)
		//	}
		//	brokerEventTypeMap[et.BrokerRef][et.GetNameAndNamespace()] = &et
		//}

		// TODO: implement triggers
		//fetchedTriggers, err := listers.TriggerLister.List(labels.Everything())
		//if err != nil {
		//	logger.Errorw("Error listing triggers", "error", err)
		//	return
		//}
		//
		//convertedTriggers := make([]Trigger, 0, len(fetchedTriggers))
		//for _, tr := range fetchedTriggers {
		//	convertedTriggers = append(convertedTriggers, convertTrigger(tr))
		//
		//	if tr.Spec.Filters != nil && len(tr.Spec.Filters) > 0 {
		//		// TODO: this is pretty hard!
		//	} else{
		//		// spec.Filter is only used when spec.Filters is nil or empty
		//		if tr.Spec.Filter != nil {
		//			for key, val := range tr.Spec.Filter.Attributes {
		//				if key == "type" {
		//					if et, ok := brokerEventTypeMap[tr.Spec.Broker][val]; ok {
		//						et.Consumers = append(et.Consumers, ObjNameAndNamespace(tr))
		//					}
		//				}
		//			}
		//		}
		//	}
		//}

		eventMesh := EventMesh{
			EventTypes: convertedEventTypes,
			Brokers:    convertedBrokers,
			// TODO
			// Triggers:   convertedTriggers,
		}

		err = json.NewEncoder(w).Encode(eventMesh)
		if err != nil {
			logger.Errorw("Error encoding event mesh", "error", err)
			return
		}
	}
}

func fetchBrokers(brokerLister eventinglistersv1.BrokerLister, logger *zap.SugaredLogger) (error, []*Broker) {
	fetchedBrokers, err := brokerLister.List(labels.Everything())
	if err != nil {
		logger.Errorw("Error listing brokers", "error", err)
		return err, nil
	}

	convertedBrokers := make([]*Broker, 0, len(fetchedBrokers))
	for _, br := range fetchedBrokers {
		convertedBroker := convertBroker(br)
		convertedBrokers = append(convertedBrokers, &convertedBroker)
	}
	return err, convertedBrokers
}
