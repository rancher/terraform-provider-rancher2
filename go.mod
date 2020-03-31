module github.com/terraform-providers/terraform-provider-rancher2

go 1.12

require (
	github.com/aws/aws-sdk-go v1.25.48 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/google/pprof v0.0.0-20190723021845-34ac40c74b70 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/terraform-plugin-sdk v1.7.0
	github.com/rancher/norman v0.0.0-20200321231028-b5f2e33b54fa
	github.com/rancher/rke v1.1.0
	github.com/rancher/types v0.0.0-20200326224903-b4612bd96d9b
	golang.org/x/crypto v0.0.0-20191112222119-e1110fd1c708
	gopkg.in/yaml.v2 v2.2.5
	k8s.io/apiserver v0.17.2
)

replace (
	k8s.io/api => k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery => github.com/rancher/apimachinery v0.17.2-rancher.1
	k8s.io/apiserver => k8s.io/apiserver v0.17.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.2
	k8s.io/client-go => github.com/rancher/client-go v1.17.2-rancher.3
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.2
	k8s.io/code-generator => k8s.io/code-generator v0.17.2
	k8s.io/component-base => k8s.io/component-base v0.17.2
	k8s.io/cri-api => k8s.io/cri-api v0.17.2
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.2
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.2
	k8s.io/kubectl => k8s.io/kubectl v0.17.2
	k8s.io/kubelet => k8s.io/kubelet v0.17.2
	k8s.io/kubernetes => k8s.io/kubernetes v1.17.2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.2
	k8s.io/metrics => k8s.io/metrics v0.17.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.2
	launchpad.net/gocheck => github.com/go-check/check v0.0.0-20200227125254-8fa46927fb4f
)
