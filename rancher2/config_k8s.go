package rancher2

import (
	"context"
	"fmt"

	lasso "github.com/rancher/lasso/pkg/client"
	lassoScheme "github.com/rancher/lasso/pkg/scheme"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	apischema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	apiV1 "k8s.io/client-go/tools/clientcmd/api/v1"
)

const (
	configMapAPIGroup   = ""
	configMapAPIVersion = "v1"
	configMapKind       = "ConfigMap"
	secretAPIGroup      = ""
	secretAPIVersion    = "v1"
	secretKind          = "Secret"
)

type k8sFactory struct {
	cli lasso.SharedClientFactory
}

func newK8sFactory(config *rest.Config) (*k8sFactory, error) {
	if config == nil {
		return nil, fmt.Errorf("Config is nil")
	}
	opts := &lasso.SharedClientFactoryOptions{
		Scheme: lassoScheme.All,
	}
	opts.Scheme.AddKnownTypeWithName(apischema.GroupVersionKind{Group: configMapAPIGroup, Version: configMapAPIVersion, Kind: configMapKind}, &corev1.ConfigMap{})
	opts.Scheme.AddKnownTypeWithName(apischema.GroupVersionKind{Group: configMapAPIGroup, Version: configMapAPIVersion, Kind: configMapKind + "List"}, &corev1.ConfigMapList{})
	opts.Scheme.AddKnownTypeWithName(apischema.GroupVersionKind{Group: secretAPIGroup, Version: secretAPIVersion, Kind: secretKind}, &corev1.Secret{})
	opts.Scheme.AddKnownTypeWithName(apischema.GroupVersionKind{Group: secretAPIGroup, Version: secretAPIVersion, Kind: secretKind + "List"}, &corev1.SecretList{})
	cli, err := lasso.NewSharedClientFactory(config, opts)
	if err != nil {
		return nil, fmt.Errorf("Error creating factory: %s", err)
	}
	return &k8sFactory{cli: cli}, err
}

func (c *k8sFactory) k8sClientForKind(gvk apischema.GroupVersionKind) (k8sClientInterface, error) {
	cli, err := c.cli.ForKind(gvk)
	if err != nil {
		return nil, err
	}
	obj, objList, err := c.cli.NewObjects(gvk)
	if err != nil {
		return nil, fmt.Errorf("Error creating new %s object: %s", gvk.Kind, err)
	}
	out := &k8sClient{
		cli:     cli,
		obj:     obj,
		objList: objList,
	}

	return out, nil
}

func (c *k8sFactory) isNamespaced(gvk apischema.GroupVersionKind) bool {
	_, namespaced, _ := c.cli.ResourceForGVK(gvk)
	return namespaced
}

func (c *k8sFactory) newCatalogV2Client() (k8sClientInterface, error) {
	return c.k8sClientForKind(apischema.GroupVersionKind{Group: catalogV2APIGroup, Version: catalogV2APIVersion, Kind: catalogV2Kind})
}

func (c *k8sFactory) newConfigMapClient() (k8sClientInterface, error) {
	return c.k8sClientForKind(apischema.GroupVersionKind{Group: configMapAPIGroup, Version: configMapAPIVersion, Kind: configMapKind})
}

func (c *k8sFactory) newSecretClient() (k8sClientInterface, error) {
	return c.k8sClientForKind(apischema.GroupVersionKind{Group: secretAPIGroup, Version: secretAPIVersion, Kind: secretKind})
}

type k8sClientInterface interface {
	Create(namespace string, obj runtime.Object) (runtime.Object, error)
	Update(namespace string, obj runtime.Object) (runtime.Object, error)
	UpdateStatus(namespace string, obj runtime.Object) (runtime.Object, error)
	Delete(name, namespace string, options *metav1.DeleteOptions) error
	Get(name, namespace string, options metav1.GetOptions) (runtime.Object, error)
	List(namespace string, opts metav1.ListOptions) (runtime.Object, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Cli() *lasso.Client
}

type k8sClient struct {
	cli     *lasso.Client
	obj     runtime.Object
	objList runtime.Object
}

func (c *k8sClient) Cli() *lasso.Client {
	return c.cli
}

func (c *k8sClient) Create(namespace string, obj runtime.Object) (runtime.Object, error) {
	result := c.obj.DeepCopyObject()
	return result, c.cli.Create(context.TODO(), namespace, obj, result, metav1.CreateOptions{})
}

func (c *k8sClient) Update(namespace string, obj runtime.Object) (runtime.Object, error) {
	result := c.obj.DeepCopyObject()
	return result, c.cli.Update(context.TODO(), namespace, obj, result, metav1.UpdateOptions{})
}

func (c *k8sClient) UpdateStatus(namespace string, obj runtime.Object) (runtime.Object, error) {
	result := c.obj.DeepCopyObject()
	return result, c.cli.UpdateStatus(context.TODO(), namespace, obj, result, metav1.UpdateOptions{})
}

func (c *k8sClient) Delete(name, namespace string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.cli.Delete(context.TODO(), namespace, name, *options)
}

func (c *k8sClient) Get(name, namespace string, options metav1.GetOptions) (runtime.Object, error) {
	result := c.obj.DeepCopyObject()
	return result, c.cli.Get(context.TODO(), namespace, name, result, options)
}

func (c *k8sClient) List(namespace string, opts metav1.ListOptions) (runtime.Object, error) {
	result := c.objList.DeepCopyObject()
	return result, c.cli.List(context.TODO(), namespace, result, opts)
}

func (c *k8sClient) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.cli.Watch(context.TODO(), namespace, opts)
}

func getK8sRestConfig(config string) (*rest.Config, error) {
	if len(config) == 0 {
		return nil, fmt.Errorf("Config is nil")
	}
	kubeJSON, err := YAMLToJSON(config)
	if err != nil {
		return nil, fmt.Errorf("Error parsing kubeconfig yaml")
	}
	kubeConfig := &apiV1.Config{}
	err = jsonToInterface(kubeJSON, kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Error parsing kubeconfig json\n%s", kubeJSON)
	}
	currentContext := kubeConfig.CurrentContext
	if len(currentContext) == 0 {
		return nil, fmt.Errorf("Current context is nil")
	}
	var context *apiV1.Context
	for i := range kubeConfig.Contexts {
		if kubeConfig.Contexts[i].Name == currentContext {
			context = &kubeConfig.Contexts[i].Context
			break
		}
	}
	if context == nil {
		return nil, fmt.Errorf("Context %s is nil", currentContext)
	}
	var auth *apiV1.AuthInfo
	for i := range kubeConfig.AuthInfos {
		if kubeConfig.AuthInfos[i].Name == context.AuthInfo {
			auth = &kubeConfig.AuthInfos[i].AuthInfo
			break
		}
	}
	if auth == nil {
		return nil, fmt.Errorf("User on context %s is nil", context)
	}
	var cluster *apiV1.Cluster
	for i := range kubeConfig.Clusters {
		if kubeConfig.Clusters[i].Name == context.Cluster {
			cluster = &kubeConfig.Clusters[i].Cluster
			break
		}
	}
	if cluster == nil {
		return nil, fmt.Errorf("Cluster on context %s is nil", context)
	}
	restConfig := &rest.Config{}
	restConfig.BearerToken = auth.Token
	restConfig.BearerTokenFile = auth.TokenFile
	restConfig.Username = auth.Username
	restConfig.Password = auth.Password
	restConfig.CertFile = auth.ClientCertificate
	restConfig.KeyFile = auth.ClientKey
	restConfig.CertData = auth.ClientCertificateData
	restConfig.KeyData = auth.ClientKeyData
	restConfig.Host = cluster.Server
	restConfig.CAFile = cluster.CertificateAuthority
	restConfig.CAData = cluster.CertificateAuthorityData
	restConfig.Insecure = cluster.InsecureSkipTLSVerify
	return restConfig, nil
}
