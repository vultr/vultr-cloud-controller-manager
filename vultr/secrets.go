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

type SecretWatch struct {
	kubeClient kubernetes.Interface
	ctx        context.Context
	Services   []*v1.Service
	secrets    map[string][]SecretList
}

type SecretList struct {
	Name    string
	Service string
}

var SecretWatcher SecretWatch

func SetupSecretWatcher(ctx context.Context) {
	SecretWatcher = SecretWatch{ctx: ctx, secrets: make(map[string][]SecretList)}
}

func (s *SecretWatch) AddService(svc *v1.Service, secretName string) {
	s.Services = append(s.Services, svc)
	// [namespace] -> ["secret-name/service-name"]
	// Example [nginx] -> ["prod-tls-cert/nginx-frontend"]
	s.secrets[svc.Namespace] = append(s.secrets[svc.Namespace], SecretList{Service: svc.Name, Name: secretName})
	klog.Infof("added secret %s to watcher", secretName)
}

func (s *SecretWatch) WatchSecrets() {
	if err := s.getKubeClient(); err != nil {
		klog.V(3).Info(err)
	}

	watcher, err := s.kubeClient.CoreV1().Secrets("*").Watch(s.ctx, metav1.ListOptions{})
	if err != nil {
		klog.V(3).Info(err)
	}

	for event := range watcher.ResultChan() {
		secret := event.Object.(*v1.Secret)

		switch event.Type {
		case watch.Modified:
			if _, ok := s.secrets[secret.Namespace]; ok {
				for _, sec := range s.secrets[secret.Namespace] {
					if sec.Name == secret.Name {
						s.updateServicefromSecret(sec.Service, secret.Namespace)
					}
				}
			}

		default:
			continue
		}

	}

}

func (s *SecretWatch) updateServicefromSecret(svcName, namespace string) {
	if err := s.getKubeClient(); err != nil {
		klog.V(3).Info(err)
	}

	svc, err := s.kubeClient.CoreV1().Services(namespace).Get(s.ctx, svcName, metav1.GetOptions{})
	if err != nil {
		klog.V(3).Info(err)
	}

	svc.Annotations[annoVultrLBSSLLastUpdatedTime] = time.Now().String()

	s.kubeClient.CoreV1().Services(namespace).Update(s.ctx, svc, metav1.UpdateOptions{})

	klog.V(3).Infof("service %s in namespace %s has been updated", svcName, namespace)
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
