package rancher2

import (
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rancher/norman/clientbase"
	"github.com/stretchr/testify/assert"
)

func testAppV2Operation(transitioning, opError bool) map[string]interface{} {
	return map[string]interface{}{
		"metadata": map[string]interface{}{
			"state": map[string]interface{}{
				"transitioning": transitioning,
				"error":         opError,
				"message":       "operation message",
			},
		},
	}
}

func testAppV2APIError(statusCode int) error {
	return &clientbase.APIError{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Msg:        http.StatusText(statusCode),
	}
}

type testAppV2OperationResult struct {
	op  map[string]interface{}
	err error
}

// testAppV2OperationGetter returns the given results in order, repeating the last one.
func testAppV2OperationGetter(results ...testAppV2OperationResult) (func() (map[string]interface{}, error), *atomic.Int32) {
	calls := &atomic.Int32{}
	return func() (map[string]interface{}, error) {
		n := int(calls.Add(1))
		r := results[min(n, len(results))-1]
		return r.op, r.err
	}, calls
}

func testAppV2OperationLogs(logs string, err error) func(map[string]interface{}) (string, error) {
	return func(map[string]interface{}) (string, error) {
		return logs, err
	}
}

func TestWaitForAppV2Operation(t *testing.T) {
	permanentErr := errors.New("permanent error")
	cases := []struct {
		name      string
		results   []testAppV2OperationResult
		logs      string
		logsErr   error
		timeout   time.Duration
		wantErr   string
		wantCalls int
	}{
		{
			name: "succeeds after transient errors",
			results: []testAppV2OperationResult{
				{err: testAppV2APIError(http.StatusNotFound)},
				{err: testAppV2APIError(http.StatusServiceUnavailable)},
				{op: testAppV2Operation(true, false)},
				{op: testAppV2Operation(false, false)},
			},
			timeout:   time.Minute,
			wantCalls: 4,
		},
		{
			name: "returns non-transient errors immediately",
			results: []testAppV2OperationResult{
				{err: permanentErr},
				{op: testAppV2Operation(false, false)},
			},
			timeout:   time.Minute,
			wantErr:   permanentErr.Error(),
			wantCalls: 1,
		},
		{
			name: "returns operation logs when the operation fails",
			results: []testAppV2OperationResult{
				{op: testAppV2Operation(false, true)},
			},
			logs:      "helm install failed",
			timeout:   time.Minute,
			wantErr:   "helm install failed",
			wantCalls: 1,
		},
		{
			name: "returns operation message when logs are unavailable",
			results: []testAppV2OperationResult{
				{op: testAppV2Operation(false, true)},
			},
			logsErr:   errors.New("logs unavailable"),
			timeout:   time.Minute,
			wantErr:   "operation message: logs unavailable",
			wantCalls: 1,
		},
		{
			name: "returns last transient error after timeout",
			results: []testAppV2OperationResult{
				{err: testAppV2APIError(http.StatusNotFound)},
			},
			timeout: 2 * time.Second,
			wantErr: http.StatusText(http.StatusNotFound),
		},
		{
			name: "times out while operation is in progress",
			results: []testAppV2OperationResult{
				{op: testAppV2Operation(true, false)},
			},
			timeout: 2 * time.Second,
			wantErr: "App V2 operation ns/op is still in progress",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			getOperation, calls := testAppV2OperationGetter(tc.results...)
			err := waitForAppV2Operation("ns/op", tc.timeout, getOperation, testAppV2OperationLogs(tc.logs, tc.logsErr))
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErr)
			}
			if tc.wantCalls > 0 {
				assert.Equal(t, int32(tc.wantCalls), calls.Load())
			} else {
				assert.Greater(t, calls.Load(), int32(1))
			}
		})
	}
}

func TestIsAppV2OperationTransientError(t *testing.T) {
	cases := []struct {
		err  error
		want bool
	}{
		{err: testAppV2APIError(http.StatusNotFound), want: true},
		{err: testAppV2APIError(http.StatusInternalServerError), want: true},
		{err: testAppV2APIError(http.StatusBadGateway), want: true},
		{err: testAppV2APIError(http.StatusServiceUnavailable), want: true},
		{err: testAppV2APIError(http.StatusForbidden), want: false},
		{err: testAppV2APIError(http.StatusUnauthorized), want: false},
		{err: errors.New("permanent error"), want: false},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.want, isAppV2OperationTransientError(tc.err), tc.err.Error())
	}
}
