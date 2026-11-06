package veille

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"
)

// JobTerminalStatuses are the job statuses WaitForJob stops on.
var JobTerminalStatuses = []string{"CompletedJob", "Error", "NotFound", "BadRequest"}

// WaitForJobOptions configures WaitForJob. The zero value polls after 3 s, then 1.5x longer each time up to
// 60 s, for up to 12 hours.
type WaitForJobOptions struct {
	Interval   time.Duration
	Timeout    time.Duration
	OnProgress func(job map[string]any)
}

// WaitForJob polls GET /v1/jobs/{id} until the job reaches a terminal status and returns the last response.
func (c *Client) WaitForJob(ctx context.Context, jobID string, options *WaitForJobOptions) (*Response, error) {
	if strings.TrimSpace(jobID) == "" {
		return nil, &SerializationError{Param: "id", Message: "A non-blank job id is required."}
	}
	opts := WaitForJobOptions{}
	if options != nil {
		opts = *options
	}
	if opts.Interval <= 0 {
		opts.Interval = 3 * time.Second
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 12 * time.Hour
	}
	started := time.Now()
	delay := opts.Interval
	maxDelay := max(delay, time.Minute)
	for {
		last, err := c.Request(ctx, http.MethodGet, "/v1/jobs/{id}", map[string]any{"id": jobID})
		if err != nil {
			return nil, err
		}
		job, err := last.Map()
		if err != nil {
			return nil, err
		}
		if opts.OnProgress != nil {
			opts.OnProgress(job)
		}
		if status, ok := job["status"].(string); ok && slices.Contains(JobTerminalStatuses, status) {
			return last, nil
		}
		if time.Since(started)+delay > opts.Timeout {
			return nil, &JobTimeoutError{JobID: jobID, LastResponse: last}
		}
		timer := time.NewTimer(delay)
		delay = min(delay*3/2, maxDelay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, &TransportError{Message: "WaitForJob was cancelled", Err: ctx.Err()}
		case <-timer.C:
		}
	}
}
