package registry

import (
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ClusterRegistry struct {
	namespace            string
	dockerjsonSecretName string
	coreClient           kubernetes.Interface
}

var _ RegistryInfoSource = ClusterRegistry{}

func NewClusterRegistry(namespace string, coreClient kubernetes.Interface) ClusterRegistry {
	if len(namespace) == 0 {
		namespace = "kectl-cluster-registry"
	}
	return ClusterRegistry{
		namespace:            namespace,
		dockerjsonSecretName: "registry-dockerjson",
		coreClient:           coreClient,
	}
}

func (r ClusterRegistry) Info() (RegistryInfo, error) {
	secretDesc := fmt.Sprintf("secret '%s' in namespace '%s'", r.dockerjsonSecretName, r.namespace)

	secret, err := r.coreClient.CoreV1().Secrets(r.namespace).Get(r.dockerjsonSecretName, metav1.GetOptions{})
	if err != nil {
		return RegistryInfo{}, fmt.Errorf("Getting dockerjson %s: %s", secretDesc, err)
	}

	const addressAnnKey = "registry-address" // TODO namespace?

	address := secret.Annotations[addressAnnKey]
	if len(address) == 0 {
		return RegistryInfo{}, fmt.Errorf("Expected to find non-empty annotation '%s' on %s", addressAnnKey, secretDesc)
	}

	dockerJSON := secret.Data[corev1.DockerConfigJsonKey]

	username, password, err := r.extractUserPass(dockerJSON)
	if err != nil {
		return RegistryInfo{}, fmt.Errorf("Parsing Docker config from %s: %s", secretDesc, err)
	}

	info := RegistryInfo{
		Address:    address,
		Username:   username,
		Password:   password,
		Dockerjson: dockerJSON,
	}

	return info, nil
}

func (r ClusterRegistry) extractUserPass(dockerJSON []byte) (string, string, error) {
	var conf dockerConfig

	err := json.Unmarshal(dockerJSON, &conf)
	if err != nil {
		return "", "", err
	}

	if len(conf.Auths) != 1 {
		return "", "", fmt.Errorf("Expected one auth configuration")
	}

	for _, v := range conf.Auths {
		username := v["username"]
		password := v["password"]

		if len(username) == 0 || len(password) == 0 {
			return "", "", fmt.Errorf("Expected non-empty username and password")
		}

		return username, password, nil
	}

	panic("Unreachable")
}

type dockerConfig struct {
	Auths dockerConfigAuths
}

type dockerConfigAuths map[string]map[string]string
