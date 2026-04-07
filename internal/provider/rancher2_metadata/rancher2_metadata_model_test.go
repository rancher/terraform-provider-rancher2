package rancher2_metadata

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestMetadataToTypesObject(t *testing.T) {
	testCases := []struct {
		name string
		fit  Metadata
		want types.Object
	}{
		{
			"Basic",
			SampleMetadataGoModel(),
			SampleMetadataTypesObject(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			diags := &diag.Diagnostics{}
			got := tc.fit.ToTypesObject(ctx, diags)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
				t.Errorf("unexpected diff (-want, +got) = %s", diff)
			}
		})
	}
}

func TestOwnerReferenceToTypesObject(t *testing.T) {
	testCases := []struct {
		name string
		fit  OwnerReference
		want types.Object
	}{
		{
			"Basic",
			OwnerReference{
				APIVersion:         "v1",
				Kind:               "some_kind",
				Name:               "owner",
				UID:                "some_uid",
				Controller:         true,
				BlockOwnerDeletion: true,
			},
			types.ObjectValueMust(
				OwnerReferenceAttrTypes,
				map[string]attr.Value{
					"api_version":          types.StringValue("v1"),
					"kind":                 types.StringValue("some_kind"),
					"name":                 types.StringValue("owner"),
					"uid":                  types.StringValue("some_uid"),
					"controller":           types.BoolValue(true),
					"block_owner_deletion": types.BoolValue(true),
				},
			),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			diags := &diag.Diagnostics{}
			got := tc.fit.ToTypesObject(ctx, diags)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if diff := cmp.Diff(tc.want, got, cmp.AllowUnexported(tftypes.Value{})); diff != "" {
				t.Errorf("unexpected diff (-want, +got) = %s", diff)
			}
		})
	}
}

func TestToGoModel(t *testing.T) {
	testCases := []struct {
		name string
		fit  types.Object
		want *Metadata
	}{
		{
			"Basic",
			SampleMetadataTypesObject(),
			func() *Metadata {
				m := SampleMetadataGoModel()
				return &m
			}(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			diags := &diag.Diagnostics{}
			got := ToGoModel(ctx, diags, tc.fit)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unexpected diff (-want, +got) = %s", diff)
			}
		})
	}
}
