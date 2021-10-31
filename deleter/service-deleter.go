package deleter

import (
	log "github.com/sirupsen/logrus"

	"github.com/cilium/cilium/pkg/client"
)

type NodePortDeleter struct {
	client *client.Client
}

func New(url string) (NodePortDeleter, error) {
	c, err := client.NewClient(url)
	if err != nil {
		return NodePortDeleter{}, err
	}

	cc := NodePortDeleter{
		client: c,
	}

	return cc, nil
}

func NewDefaultclient() (NodePortDeleter, error) {
	c, err := client.NewDefaultClient()
	if err != nil {
		return NodePortDeleter{}, err
	}

	cc := NodePortDeleter{
		client: c,
	}

	return cc, nil
}

func (cc *NodePortDeleter) DeleteServices() error {
	for {
		services, err := cc.client.GetServices()
		if err != nil {
			return err
		}

		for _, service := range services {
			// From docs: https://pkg.go.dev/github.com/cilium/cilium@v1.10.0/api/v1/models#ServiceSpecFlags
			// type ServiceSpecFlags struct {
			//     // Service type
			//     // Enum: [ClusterIP NodePort ExternalIPs HostPort LoadBalancer LocalRedirect]
			//     Type string `json:"type,omitempty"`
			// }

			if service.Spec.Flags.Type == "ClusterIP" || service.Spec.Flags.Type == "LoadBalancer" {
				log.WithFields(log.Fields{
					"ns":   service.Spec.Flags.Namespace,
					"name": service.Spec.Flags.Name,
				}).Debugf("ignored %s", service.Spec.Flags.Type)
				continue
			}

			if service.Spec.Flags.Type == "HostPort" {
				log.WithFields(log.Fields{
					"ns":   service.Spec.Flags.Namespace,
					"name": service.Spec.Flags.Name,
				}).Warningf("A hostport is open!")
			}

			if service.Spec.Flags.Type == "ExternalIPs" {
				log.WithFields(log.Fields{
					"ns":   service.Spec.Flags.Namespace,
					"name": service.Spec.Flags.Name,
				}).Warningf("ExternalIPs are being used. change it to LoadBalancer with a proper MetalLB configuration.")
			}

			log.WithFields(log.Fields{
				"ns":   service.Spec.Flags.Namespace,
				"name": service.Spec.Flags.Name,
				"type": service.Spec.Flags.Type,
			}).Debugf("deleting service")

			if err := cc.client.DeleteServiceID(service.Spec.ID); err != nil {
				log.WithFields(log.Fields{
					"ns":   service.Spec.Flags.Namespace,
					"name": service.Spec.Flags.Name,
					"type": service.Spec.Flags.Type,
					"ID":   service.Spec.ID,
				}).Errorf("error trying to delete service")
				return err
			} else {
				log.WithFields(log.Fields{
					"ns":   service.Spec.Flags.Namespace,
					"name": service.Spec.Flags.Name,
					"type": service.Spec.Flags.Type,
					"ID":   service.Spec.ID,
				}).Infof("service deleted")
			}
		}

		return nil
	}
}
