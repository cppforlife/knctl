package registry

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ImageInfo struct {
	Address            string
	ServiceAccountName string
}

type RegistryInfo struct {
	Address    string
	Username   string
	Password   string
	Dockerjson []byte
}

type RegistryInfoSource interface {
	Info() (RegistryInfo, error)
}

type Registry struct {
	nsForImagePull      string
	imagePullSecretName string
	imagePushSecretName string
	serviceAccountName  string
	coreClient          kubernetes.Interface
	registryInfoSrc     RegistryInfoSource
}

func NewRegistry(nsForImagePull string, coreClient kubernetes.Interface, registryInfoSrc RegistryInfoSource) Registry {
	return Registry{
		nsForImagePull:      nsForImagePull,
		imagePullSecretName: "knctl-pull-cluster-registry",
		imagePushSecretName: "knctl-push-cluster-registry",
		serviceAccountName:  "knctl-cluster-registry",
		coreClient:          coreClient,
		registryInfoSrc:     registryInfoSrc,
	}
}

func (r Registry) Info() (ImageInfo, error) {
	var info ImageInfo

	registryInfo, err := r.registryInfoSrc.Info()
	if err != nil {
		return info, fmt.Errorf("Getting registry info: %s", err)
	}

	err = r.createOrUpdateImagePullSecret(registryInfo)
	if err != nil {
		return info, err
	}

	err = r.createOrUpdateImagePushSecret(registryInfo)
	if err != nil {
		return info, err
	}

	err = r.createOrUpdateServiceAccount()
	if err != nil {
		return info, err
	}

	info = ImageInfo{
		Address:            registryInfo.Address,
		ServiceAccountName: r.serviceAccountName,
	}

	return info, nil
}

func (r Registry) createOrUpdateImagePullSecret(registryInfo RegistryInfo) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.imagePullSecretName,
			Namespace: r.nsForImagePull,
		},
		Type: corev1.SecretTypeDockerConfigJson,
		StringData: map[string]string{
			corev1.DockerConfigJsonKey: string(registryInfo.Dockerjson),
		},
	}

	_, err := r.coreClient.CoreV1().Secrets(r.nsForImagePull).Create(secret)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return r.updateSecret(secret)
		}
		return fmt.Errorf("Creating image pull secret: %s", err)
	}

	return nil
}

func (r Registry) createOrUpdateImagePushSecret(registryInfo RegistryInfo) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.imagePushSecretName,
			Namespace: r.nsForImagePull,
			Annotations: map[string]string{
				"build.knative.dev/docker-0": "https://" + registryInfo.Address,
			},
		},
		Type: corev1.SecretTypeBasicAuth,
		StringData: map[string]string{
			corev1.BasicAuthUsernameKey: registryInfo.Username,
			corev1.BasicAuthPasswordKey: registryInfo.Password,
		},
	}

	_, err := r.coreClient.CoreV1().Secrets(r.nsForImagePull).Create(secret)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return r.updateSecret(secret)
		}
		return fmt.Errorf("Creating image push secret: %s", err)
	}

	return nil
}

func (r Registry) updateSecret(secret *corev1.Secret) error {
	foundSecret, err := r.coreClient.CoreV1().Secrets(secret.Namespace).Get(secret.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	foundSecret.Annotations = secret.Annotations
	foundSecret.Type = secret.Type
	foundSecret.StringData = secret.StringData

	_, err = r.coreClient.CoreV1().Secrets(foundSecret.Namespace).Update(foundSecret)
	if err != nil {
		return fmt.Errorf("Updating secret: %s", err)
	}

	return nil
}

func (r Registry) createOrUpdateServiceAccount() error {
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.serviceAccountName,
			Namespace: r.nsForImagePull,
		},
		Secrets:          []corev1.ObjectReference{{Name: r.imagePushSecretName}},
		ImagePullSecrets: []corev1.LocalObjectReference{{Name: r.imagePullSecretName}},
	}

	_, err := r.coreClient.CoreV1().ServiceAccounts(r.nsForImagePull).Create(serviceAccount)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return r.updateServiceAccount(serviceAccount)
		}
		return fmt.Errorf("Creating service account: %s", err)
	}

	return nil
}

func (r Registry) updateServiceAccount(serviceAccount *corev1.ServiceAccount) error {
	foundSA, err := r.coreClient.CoreV1().ServiceAccounts(serviceAccount.Namespace).Get(serviceAccount.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	foundSA.Secrets = serviceAccount.Secrets
	foundSA.ImagePullSecrets = serviceAccount.ImagePullSecrets

	_, err = r.coreClient.CoreV1().ServiceAccounts(foundSA.Namespace).Update(foundSA)
	if err != nil {
		return fmt.Errorf("Updating service account: %s", err)
	}

	return nil
}
