module github.com/rancher/terraform-provider-rancher2

go 1.16

require (
	github.com/ghodss/yaml v1.0.0
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/terraform-plugin-sdk v1.14.0
	github.com/rancher/norman v0.0.0-20211201154850-abe17976423e
	github.com/rancher/rancher v0.0.0-20211218171736-3c1d5fac3248
	github.com/rancher/rancher/pkg/apis v0.0.0
	github.com/rancher/rancher/pkg/client v0.0.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.22.3
	k8s.io/apimachinery v0.22.3
	k8s.io/apiserver v0.22.3
	k8s.io/client-go v12.0.0+incompatible
)

replace (
	github.com/crewjam/saml => github.com/crewjam/saml v0.4.1
	github.com/docker/distribution => github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker => github.com/docker/docker v20.10.6+incompatible
	github.com/rancher/rancher/pkg/apis => github.com/rancher/rancher/pkg/apis v0.0.0-20211218171736-3c1d5fac3248
	github.com/rancher/rancher/pkg/client => github.com/rancher/rancher/pkg/client v0.0.0-20211218171736-3c1d5fac3248
	github.com/spf13/afero => github.com/spf13/afero v1.2.2
	helm.sh/helm/v3 => github.com/rancher/helm/v3 v3.5.4-rancher.1
	k8s.io/api => k8s.io/api v0.22.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.22.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.3
	k8s.io/apiserver => k8s.io/apiserver v0.22.3
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.22.3
	k8s.io/client-go => github.com/rancher/client-go v1.22.3-rancher.1
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.22.3
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.22.3
	k8s.io/code-generator => k8s.io/code-generator v0.22.3
	k8s.io/component-base => k8s.io/component-base v0.22.3
	k8s.io/component-helpers => k8s.io/component-helpers v0.22.3
	k8s.io/controller-manager => k8s.io/controller-manager v0.22.3
	k8s.io/cri-api => k8s.io/cri-api v0.22.3
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.22.3
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.22.3
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.22.3
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.22.3
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.22.3
	k8s.io/kubectl => k8s.io/kubectl v0.22.3
	k8s.io/kubelet => k8s.io/kubelet v0.22.3
	k8s.io/kubernetes => k8s.io/kubernetes v1.22.3
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.22.3
	k8s.io/metrics => k8s.io/metrics v0.22.3
	k8s.io/mount-utils => k8s.io/mount-utils v0.22.3
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.22.0
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.22.3
	launchpad.net/gocheck => github.com/go-check/check v0.0.0-20200227125254-8fa46927fb4f
)
