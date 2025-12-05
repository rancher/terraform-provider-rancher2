module github.com/rancher/terraform-provider-rancher2

go 1.25.4

require (
	github.com/hashicorp/terraform-plugin-framework v1.16.1
	github.com/hashicorp/terraform-plugin-go v0.29.0
	github.com/hashicorp/terraform-plugin-log v0.9.0
)

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-plugin v1.7.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/terraform-registry-address v0.4.0 // indirect
	github.com/hashicorp/terraform-svchost v0.1.1 // indirect
	github.com/hashicorp/yamux v0.1.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/grpc v1.75.1 // indirect
	google.golang.org/protobuf v1.36.9 // indirect
)

replace (
	github.com/crewjam/saml => github.com/crewjam/saml v0.4.5
	github.com/docker/distribution => github.com/docker/distribution v2.8.1+incompatible
	github.com/docker/docker => github.com/docker/docker v20.10.17+incompatible
	github.com/rancher/rancher/pkg/apis => github.com/rancher/rancher/pkg/apis v0.0.0-20250227174106-7829dbe62d7f
	github.com/rancher/rancher/pkg/client => github.com/rancher/rancher/pkg/client v0.0.0-20250227174106-7829dbe62d7f
	github.com/spf13/afero => github.com/spf13/afero v1.2.2
	go.qase.io/client => github.com/rancher/qase-go/client v0.0.0-20231114201952-65195ec001fa
	helm.sh/helm/v3 => github.com/rancher/helm/v3 v3.9.0-rancher1
	k8s.io/api => k8s.io/api v0.32.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.32.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.32.2
	k8s.io/apiserver => k8s.io/apiserver v0.32.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.32.2
	k8s.io/client-go => k8s.io/client-go v0.32.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.32.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.32.2
	k8s.io/code-generator => k8s.io/code-generator v0.32.2
	k8s.io/component-base => k8s.io/component-base v0.32.2
	k8s.io/component-helpers => k8s.io/component-helpers v0.32.2
	k8s.io/controller-manager => k8s.io/controller-manager v0.32.2
	k8s.io/cri-api => k8s.io/cri-api v0.32.2
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.32.2
	k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.32.2
	k8s.io/endpointslice => k8s.io/endpointslice v0.32.2
	k8s.io/kms => k8s.io/kms v0.32.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.32.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.32.2
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20230717233707-2695361300d9
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.32.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.32.2
	k8s.io/kubectl => k8s.io/kubectl v0.32.2
	k8s.io/kubelet => k8s.io/kubelet v0.32.2
	k8s.io/kubernetes => k8s.io/kubernetes v1.30.1
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.32.2
	k8s.io/metrics => k8s.io/metrics v0.32.2
	k8s.io/mount-utils => k8s.io/mount-utils v0.32.2
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.32.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.32.2
	launchpad.net/gocheck => github.com/go-check/check v0.0.0-20200227125254-8fa46927fb4f
)
