package rancher2_dev2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
)

type Rancher2Dev2ResourceModel struct {
	ID         types.String `tfsdk:"id"`
	APIVersion types.String `tfsdk:"api_version"`
	Kind       types.String `tfsdk:"kind"`
	Metadata   types.Object `tfsdk:"metadata"`
	Spec       types.Object `tfsdk:"spec"`
	Status     types.String `tfsdk:"status"`
}


func (m *Rancher2Dev2ResourceModel) ToPlan(ctx context.Context, diags *diag.Diagnostics) tfsdk.Plan {
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2ResourceModel to tfsdk.Plan: \n%+v", pp.PrettyPrint(m)))
	if diags.HasError() {
		return tfsdk.Plan{}
	}
	r := NewRancher2Dev2Resource()
	s := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, s)

	plan := tfsdk.Plan{
		Schema: s.Schema,
	}
	if diags.HasError() {
		return plan
	}
	diags.Append(plan.Set(ctx, m)...)
	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2ResourceModel to tfsdk.Plan: \n%+v", pp.PrettyPrint(plan)))
	return plan
}

func (m *Rancher2Dev2ResourceModel) ToState(ctx context.Context, diags *diag.Diagnostics) tfsdk.State {
	tflog.Debug(ctx, fmt.Sprintf("Converting Rancher2Dev2ResourceModel to tfsdk.State: \n%+v", pp.PrettyPrint(m)))
	if diags.HasError() {
		return tfsdk.State{}
	}
	r := NewRancher2Dev2Resource()
	s := &resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, s)

	state := tfsdk.State{
		Schema: s.Schema,
	}
	diags.Append(state.Set(ctx, m)...)
	tflog.Debug(ctx, fmt.Sprintf("Converted Rancher2Dev2ResourceModel to tfsdk.State: \n%+v", pp.PrettyPrint(state)))
	return state
}

