package sbmgmt

import (
	"errors"
	"fmt"
	"os"

	servicebus "github.com/Azure/azure-service-bus-go"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

// GetServiceBusNamespace finds the ASB namespace and returns it
func GetServiceBusNamespace(name string) (*servicebus.Namespace, error) {
	if err := godotenv.Load(); err != nil {
		log.Error("Failed to do godotenv.Load()")
		fmt.Println(err)
		return nil, err
	}
	connStr := os.Getenv(name)
	if connStr == "" {
		log.Error("connStr is empty")
		return nil, errors.New("Connection String is empty")
	}

	ns, err := servicebus.NewNamespace(servicebus.NamespaceWithConnectionString(connStr))
	if err != nil {
		log.Error("Failed to create the Namespace client")
		return nil, err
	}
	return ns, nil
}
