package curl

import (
	"bytes"
	"fmt"
	"lazycurl/internal/model"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Executor handles executing curl commands.
type Executor struct{}

// NewExecutor creates a new curl executor.
func NewExecutor() *Executor {
	return &Executor{}
}

// Execute runs a curl command based on the request model.
func (e *Executor) Execute(req model.Request) model.Response {
	// Build curl arguments
	args := []string{"-s", "-w", "%{http_code}\n%{time_total}", "-X", req.Method}

	// Add headers
	for k, v := range req.Headers {
		args = append(args, "-H", fmt.Sprintf("%s: %s", k, v))
	}

	// Add body if present
	if req.Body != "" {
		args = append(args, "-d", req.Body)
	}

	// Add URL
	args = append(args, req.URL)

	cmd := exec.Command("curl", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	if err != nil {
		return model.Response{
			Error:     fmt.Errorf("curl execution failed: %v\nstderr: %s", err, stderr.String()),
			TimeTaken: duration,
		}
	}

	// Parse output. We used a separator in -w to distinguish body from metadata.

	// Re-defining the command with separator
	args = []string{"-s", "-w", "\n_____LAZYCURL_METADATA_____\n%{http_code}\n%{time_total}", "-X", req.Method}
	for k, v := range req.Headers {
		args = append(args, "-H", fmt.Sprintf("%s: %s", k, v))
	}
	if req.Body != "" {
		args = append(args, "-d", req.Body)
	}
	args = append(args, req.URL)

	cmd = exec.Command("curl", args...)
	stdout.Reset()
	stderr.Reset()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return model.Response{
			Error:     fmt.Errorf("curl error: %v, stderr: %s", err, stderr.String()),
			TimeTaken: duration,
		}
	}

	fullOutput := stdout.String()
	parts := strings.Split(fullOutput, "_____LAZYCURL_METADATA_____")

	if len(parts) < 2 {
		return model.Response{
			Body:      fullOutput,
			Error:     fmt.Errorf("failed to parse curl metadata"),
			TimeTaken: duration, // Fallback to wall clock
		}
	}

	body := parts[0]
	meta := strings.TrimSpace(parts[1])
	metaLines := strings.Split(meta, "\n")

	statusCode := 0
	if len(metaLines) >= 1 {
		statusCode, _ = strconv.Atoi(strings.TrimSpace(metaLines[0]))
	}

	// curl time_total is in seconds (e.g. 0.123)
	// We'll use our own wall clock `duration` for simplicity and consistency in MVP
	// unless we really need the curl internal time.

	return model.Response{
		StatusCode: statusCode,
		Body:       body,
		Headers:    nil, // TODO: Parse headers with -D or -i if needed later
		TimeTaken:  duration,
	}
}
