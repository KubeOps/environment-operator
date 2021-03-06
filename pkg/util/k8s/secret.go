package k8s

import (
	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

// Secret is a client for interacting with secrets
type Secret struct {
	kubernetes.Interface
	Namespace string
}

// List returns the list of k8s secrets maintained by pipeline for provided client
func (client *Secret) List() ([]v1.Secret, error) {
	list, err := client.Core().Secrets(client.Namespace).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// Check if the secret exists in the namespace
func (client *Secret) Exists(secretname string) bool {

	secrets, err := client.List()
	if err != nil {
		log.Error(err.Error())
	}
	found := false
	for _, sec := range secrets {
		if sec.ObjectMeta.Name == secretname {
			found = true
		}
	}
	return found
}

// Apply updates or creates secrets in k8s
func (client *Secret) Apply(resource *v1.Secret) error {
	if client.Exists(resource.Name) {
		return client.Update(resource)
	}
	return client.Create(resource)
}

// Create creates new secret in k8s
func (client *Secret) Create(resource *v1.Secret) error {
	_, err := client.
		Core().
		Secrets(client.Namespace).
		Create(resource)
	return err
}

// Update updates existing secrets in k8s
func (client *Secret) Update(resource *v1.Secret) error {
	current, err := client.Get(resource.Name)
	if err != nil {
		return err
	}
	resource.ResourceVersion = current.GetResourceVersion()

	_, err = client.
		Core().
		Secrets(client.Namespace).
		Update(resource)
	return err
}

// Get returns secret object from the k8s by name
func (client *Secret) Get(name string) (*v1.Secret, error) {
	return client.Core().Secrets(client.Namespace).Get(name)
}
