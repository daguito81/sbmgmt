package sbmgmt

import (
	"context"
	"log"
	"time"

	servicebus "github.com/Azure/azure-service-bus-go"
)

const (
	lockDuration time.Duration = 5 * time.Second
	msgDuration  time.Duration = 10 * time.Minute
)

// GetOrBuildTopic creates a topic and returns the client or error
func GetOrBuildTopic(nsname string, topicName string) (*servicebus.Topic, *servicebus.TopicEntity, error) {

	ns, err := GetServiceBusNamespace(nsname)
	if err != nil {
		log.Fatalln("Error connecting to Service Bus: ", err)
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new Queue Manager
	tm := ns.NewTopicManager()
	te, err := tm.Get(ctx, topicName)
	if err != nil && !servicebus.IsErrNotFound(err) {
		return nil, nil, err
	}

	if te == nil {
		_, err := tm.Put(ctx, topicName)
		if err != nil {
			return nil, nil, err
		}
	}

	t, err := ns.NewTopic(topicName)
	if err != nil {
		return nil, nil, err
	}
	return t, te, nil
}

// DeleteTopic deletes the named queue from Service Bus
func DeleteTopic(nsname string, topicName string) error {
	ns, err := GetServiceBusNamespace(nsname)
	if err != nil {
		log.Fatalln("Error connecting to Service Bus: ", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete Queue
	tm := ns.NewTopicManager()
	if err = tm.Delete(ctx, topicName); err != nil {
		log.Fatalln("Error deleting Queue: ", err)
		return err
	}
	return nil

}

// GetOrBuildSubscription creates a new subscription or gets an existint subscription
func GetOrBuildSubscription(nsname string, subName string, topicName string, lockDur time.Duration, msgDur time.Duration) (*servicebus.Subscription, *servicebus.SubscriptionEntity, error) {
	ns, err := GetServiceBusNamespace(nsname)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, _, err = GetOrBuildTopic(nsname, topicName)
	if err != nil {
		log.Fatalln(err)
	}
	// Create a new Topic Manager
	sm, err := ns.NewSubscriptionManager(topicName)
	if err != nil {
		log.Fatalln(err)
	}
	se, err := sm.Get(ctx, subName)
	if err != nil && !servicebus.IsErrNotFound(err) {
		return nil, nil, err
	}
	// In case of empty, create subscription

	if se == nil {
		_, err := sm.Put(ctx,
			subName,
			servicebus.SubscriptionWithLockDuration(&lockDur),
			servicebus.SubscriptionWithMessageTimeToLive(&msgDur),
		)
		if err != nil {
			return nil, nil, err
		}
	}
	// Create sub client
	s, err := sm.Topic.NewSubscription(subName)
	if err != nil {
		return nil, nil, err
	}
	return s, se, nil
}
