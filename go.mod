module github.com/terraform-providers/terraform-provider-rancher2

go 1.12

require (
	github.com/hashicorp/go-version v1.2.0
	github.com/hashicorp/terraform-plugin-sdk v1.0.0
	github.com/rancher/norman v0.0.0-20191003174345-0ac7dd6ccb36
	github.com/rancher/types v0.0.0-20191017022642-c4c5ca1581f9
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586
)

replace k8s.io/client-go => k8s.io/client-go v12.0.0+incompatible
