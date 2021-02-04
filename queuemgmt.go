package sbmgmt

import (
	"context"
	"log"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

// GetOrBuildQueue creates a queue and returns the client or error
func GetOrBuildQueue(queueName string) (*servicebus.Queue, error) {

	ns, err := GetServiceBusNamespace()
	if err != nil {
		log.Fatalln("Error connecting to Service Bus: ", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new Queue
	qm := ns.NewQueueManager()
	qe, err := qm.Get(ctx, queueName)
	if err != nil && !servicebus.IsErrNotFound(err) {
		return nil, err
	}

	if qe == nil {
		_, err := qm.Put(ctx, queueName)
		if err != nil {
			return nil, err
		}
	}

	q, err := ns.NewQueue(queueName)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// DeleteQueue deletes the named queue from Service Bus
func DeleteQueue(queueName string) error {
	ns, err := GetServiceBusNamespace()
	if err != nil {
		log.Fatalln("Error connecting to Service Bus: ", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete Queue
	qm := ns.NewQueueManager()
	if err = qm.Delete(ctx, queueName); err != nil {
		log.Fatalln("Error deleting Queue: ", err)
		return err
	}
	return nil

}
