package test_helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
)

func GenerateTestContext(t *testing.T, buf *bytes.Buffer, ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return tflogtest.RootLogger(ctx, buf)
}

var logLevels = map[string]int{
	"TRACE": 0,
	"DEBUG": 1,
	"INFO":  2,
	"WARN":  3,
	"ERROR": 4,
}

func PrintLog(t *testing.T, buf *bytes.Buffer, level ...string) {
	minLevel := -1
	if len(level) > 0 && level[0] != "" {
		var ok bool
		minLevel, ok = logLevels[strings.ToUpper(level[0])]
		if !ok {
			t.Logf("invalid log level: %s", level[0])
			return
		}
	}

	for line := range strings.SplitSeq(buf.String(), "\n") {
		if line == "" {
			continue
		}
		var logEntry map[string]any
		err := json.Unmarshal([]byte(line), &logEntry)
		if err != nil {
			t.Logf("failed to unmarshal log line: %s", err)
			continue
		}

		if minLevel != -1 {
			logLevelStr, ok := logEntry["@level"].(string)
			if !ok {
				continue // Skip logs without a level if filtering
			}

			logLevel, ok := logLevels[strings.ToUpper(logLevelStr)]
			if !ok {
				continue // Skip unknown log levels if filtering
			}

			// fmt.Printf("Log level detected: %d\n", logLevel)
			// fmt.Printf("Min level set: %d\n", minLevel)

			if logLevel < minLevel {
				continue
			}
		}

		if msg, ok := logEntry["@message"]; ok {
			t.Log(msg)
		}
	}
}

// rsc is a pointer to the resource you want to be configured, eg. *RancherDevResource
// client is the data sent to the resource, this is what the provider configure outputs
// this is where you can inject a custom client for testing.
func GetConfiguredResource(ctx context.Context, t *testing.T, rsc resource.ResourceWithConfigure, client c.Client) error {
	req := resource.ConfigureRequest{
		ProviderData: client,
	}
	res := resource.ConfigureResponse{}
	rsc.Configure(ctx, req, &res)
	if res.Diagnostics.HasError() {
		return fmt.Errorf("error configuring resource: %s", res.Diagnostics.Errors())
	}
	return nil
}

func GetTestClient(t *testing.T, ctx context.Context) *c.TestClient {
	return c.NewTestClient(ctx, "https://rancher.example.com", "", false, false, 30, 10, &c.TokenStore{})
}

//nolint:unparam
func Cntxt(t *testing.T, logLevel string) (context.Context, func()) {
	var buf bytes.Buffer
	return GenerateTestContext(t, &buf, nil), func() { PrintLog(t, &buf, logLevel) }
}

// rsc must be an empty resource with the schema function implemented.
func Schema(ctx context.Context, rsc resource.Resource) schema.Schema {
	schemaResponseContainer := &resource.SchemaResponse{}
	rsc.Schema(ctx, resource.SchemaRequest{}, schemaResponseContainer)
	return schemaResponseContainer.Schema
}

// rsc must be an empty resource with the schema function implemented.
func State(ctx context.Context, rsc resource.Resource, model any) (tfsdk.State, error) {
	state := tfsdk.State{
		Schema: Schema(ctx, rsc),
	}
	dgs := state.Set(ctx, model)
	if dgs.HasError() {
		return state, fmt.Errorf("error setting state: %s", dgs.Errors())
	}
	return state, nil
}

// rsc must be an empty resource with the schema function implemented.
func Plan(ctx context.Context, rsc resource.Resource, model any) (tfsdk.Plan, error) {
	plan := tfsdk.Plan{
		Schema: Schema(ctx, rsc),
	}
	dgs := plan.Set(ctx, model)
	if dgs.HasError() {
		return plan, fmt.Errorf("error setting plan: %s", dgs.Errors())
	}
	return plan, nil
}
