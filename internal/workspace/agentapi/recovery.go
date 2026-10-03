package agentapi

// A status alone cannot prove a failed mutation did not take effect.
type Recovery struct {
	Category   string `json:"category"`
	NextAction string `json:"next_action"`
	Retryable  bool   `json:"retryable"`
}

func RecoveryForStatus(status int) Recovery {
	r := Recovery{Category: "operation_failed", NextAction: "Inspect the current workspace and operation status before trying again."}
	switch {
	case status == 401 || status == 403 || status == 404:
		r.Category = "access_or_availability"
		r.NextAction = "Return to accessible projects or ask your workspace administrator for help."
	case status == 409 || status == 412:
		r.Category = "conflict"
		r.NextAction = "Reload the current state and review the retained draft or scope before submitting again."
	case status == 400 || status == 422:
		r.Category = "validation"
		r.NextAction = "Correct the indicated input and review the scope before submitting again."
	case status == 429:
		r.Category = "capacity"
		r.NextAction = "Check running work and the retry delay. Inspect whether a prior mutation was accepted before resubmitting."
	case status >= 500:
		r.Category = "environment_or_provider"
		r.NextAction = "Check connection, provider configuration or runtime diagnostics. Inspect persisted work before retrying a mutation."
	}
	return r
}
