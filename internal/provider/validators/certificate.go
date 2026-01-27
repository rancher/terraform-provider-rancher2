package validators

import (
	"context"
	"encoding/pem"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &certificateValidator{}

type certificateValidator struct{}

func (v certificateValidator) Description(ctx context.Context) string {
	return "must be a valid PEM-encoded certificate"
}

func (v certificateValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v certificateValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	block, _ := pem.Decode([]byte(value))
	if block == nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid PEM Encoded Certificate",
			"The value is not a valid PEM-encoded block.",
		)
		return
	}

	if block.Type != "CERTIFICATE" {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid PEM Encoded Certificate",
			fmt.Sprintf("The PEM block type is '%s', expected 'CERTIFICATE'.", block.Type),
		)
		return
	}
}

func IsCertificate() validator.String {
	return &certificateValidator{}
}
