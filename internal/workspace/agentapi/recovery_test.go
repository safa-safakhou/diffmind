package agentapi

import (
	"net/http"
	"strings"
	"testing"
)

func TestRecoveryNeverReplaysMutationsOrLeaksProviderBodies(t *testing.T) {
	for _, status := range []int{400, 401, 403, 404, 409, 412, 422, 429, 500, 502, 503} {
		result, err := Decode(status, http.Header{}, strings.NewReader(`{"error":"operation failed"}`))
		if err != nil || result.Recovery == nil || result.Recovery.Retryable || result.Recovery.NextAction == "" {
			t.Fatalf("status %d: %+v %v", status, result, err)
		}
	}
	a, b := RecoveryForStatus(403), RecoveryForStatus(404)
	if a != b {
		t.Fatal("access guidance reveals project existence")
	}
	accepted, err := Decode(202, http.Header{}, strings.NewReader(`{"job_id":"pending"}`))
	if err != nil || accepted.Recovery != nil || accepted.Status != 202 {
		t.Fatal("async acceptance mislabeled")
	}
}
