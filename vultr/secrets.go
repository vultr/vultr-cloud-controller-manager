package vultr

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"time"
)

// SecretWatch is the main structure for the secret watcher
type SecretWatch struct {
	kubeClient kubernetes.Interface
	ctx        context.Context
	secrets    map[string][]SecretList
}

// SecretList is meant to be stored as a slice of type SecretList which stores the name of the secret and it's service
type SecretList struct {
	Name    string
	Service string
}

const (
	logLevel = 3
)

// SecretWatcher is a global variable of type SecretWatch. We use a global variable so that the SecretWatcher can be accessed globally
// The watcher is meant to be ran as a go routine and in the current CCM we would not be able to run and access it globally otherwise
var SecretWatcher SecretWatch

// SetupSecretWatcher initializes the watcher
func SetupSecretWatcher(ctx context.Context) {
	SecretWatcher = SecretWatch{ctx: ctx, secrets: make(map[string][]SecretList)}
}

// AddService adds a service to watch the corresponding secret for to the secretwatcher
func (s *SecretWatch) AddService(svc *v1.Service, secretName string) {
	// [namespace] -> ["secret-name/service-name"]
	// Example [nginx] -> ["prod-tls-cert/nginx-frontend"]
	if _, ok := s.secrets[svc.Namespace]; ok {
		for _, val := range s.secrets[svc.Namespace] {
			if val.Service == svc.Name {
				klog.Infof("service %s already exists in secret watcher, returning", svc.Name)
				return
			}
		}
	}

	s.secrets[svc.Namespace] = append(s.secrets[svc.Namespace], SecretList{Service: svc.Name, Name: secretName})
	klog.Infof("added secret %s to watcher", secretName)
}

// WatchSecrets is the main entrance into the execution of the secretwatcher
func (s *SecretWatch) WatchSecrets() {
	if err := s.getKubeClient(); err != nil {
		klog.V(logLevel).Info(err)
		return
	}

	watcher, err := s.kubeClient.CoreV1().Secrets("").Watch(s.ctx, metav1.ListOptions{Watch: true})
	if err != nil {
		klog.V(logLevel).Info(err)
		return
	}

	for event := range watcher.ResultChan() {
		secret := event.Object.(*v1.Secret)

		switch event.Type {
		case watch.Modified:
			fallthrough
		case watch.Added:
			if _, ok := s.secrets[secret.Namespace]; ok {
				for _, sec := range s.secrets[secret.Namespace] {
					if sec.Name == secret.Name {
						klog.V(logLevel).Infof("secret %s had a %s event", secret.Name, event.Type)
						s.updateServiceFromSecret(sec.Service, secret.Namespace)
					}
				}
			}
		default:
			continue
		}
	}
}

func (s *SecretWatch) updateServiceFromSecret(svcName, namespace string) {
	if err := s.getKubeClient(); err != nil {
		klog.V(logLevel).Info(err)
	}

	svc, err := s.kubeClient.CoreV1().Services(namespace).Get(s.ctx, svcName, metav1.GetOptions{})
	if err != nil {
		klog.V(logLevel).Info(err)
	}

	svc.Annotations[annoVultrLBSSLLastUpdatedTime] = time.Now().String()

	_, err = s.kubeClient.CoreV1().Services(namespace).Update(s.ctx, svc, metav1.UpdateOptions{})
	if err != nil {
		klog.V(logLevel).Info(err)
	}

	klog.V(logLevel).Infof("service %s in namespace %s has been updated", svcName, namespace)
}

func (s *SecretWatch) getKubeClient() error {
	if s.kubeClient != nil {
		return nil
	}

	var (
		kubeConfig *rest.Config
		err        error
		config     string
	)

	// If no kubeconfig was passed in or set then we want to default to an empty string
	// This will have `clientcmd.BuildConfigFromFlags` default to `restclient.InClusterConfig()` which was existing behavior
	if Options.KubeconfigFlag == nil || Options.KubeconfigFlag.Value.String() == "" {
		config = ""
	} else {
		config = Options.KubeconfigFlag.Value.String()
	}

	kubeConfig, err = clientcmd.BuildConfigFromFlags("", config)
	if err != nil {
		return err
	}

	s.kubeClient, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	return nil
}
