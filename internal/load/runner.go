package load

import (
	"lazycurl/internal/curl"
	"lazycurl/internal/model"
	"time"
)

// Runner orchestrates the load test.
type Runner struct {
	Stats    *Stats
	Executor *curl.Executor
}

// NewRunner creates a new runner.
func NewRunner() *Runner {
	return &Runner{
		Stats:    NewStats(),
		Executor: curl.NewExecutor(),
	}
}

// StatsMsg is the message sent tea-style to the UI.
type StatsMsg struct {
	Stats *Stats
	Done  bool
}

// Run executes the load test. It returns a channel that emits updates.
// We return a tea.Cmd wrapper instead? No, we need it to stream.
// In Bubbletea, we can return a "Command" that wraps a subscription, or use a channel command.
// We'll use the "Program.Program" style or "Sub" style.
// But for simplicity here, this method spawns the work and the caller (Update) will listen to a channel.
// Actually, standard Bubbletea pattern for streaming is:
// Update() returns a Cmd that waits for the next value.
// So we should return a channel, and have a "WaitForStats(chan)" command.

func (r *Runner) Run(req model.Request, concurrency int, duration time.Duration) chan StatsMsg {
	ch := make(chan StatsMsg)

	go func() {
		// Control channels
		stop := time.After(duration)
		results := make(chan model.Response, concurrency)

		// Spawn workers
		activeWorkers := 0
		for i := 0; i < concurrency; i++ {
			activeWorkers++
			go func() {
				for {
					// We need a way to stop workers cleanly.
					// For MVP, simplistic check?
					resp := r.Executor.Execute(req)
					results <- resp
				}
			}()
		}

		// Collect results
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		startTime := time.Now()

	loop:
		for {
			select {
			case <-stop:
				break loop
			case resp := <-results:
				r.updateStats(resp)
			case <-ticker.C:
				r.Stats.ElapsedTime = time.Since(startTime)
				ch <- StatsMsg{Stats: r.Stats, Done: false}
			}
		}

		// Final stats
		r.Stats.ElapsedTime = time.Since(startTime)
		ch <- StatsMsg{Stats: r.Stats, Done: true}
		close(ch)
	}()

	return ch
}

func (r *Runner) updateStats(resp model.Response) {
	r.Stats.TotalRequests++
	r.Stats.StatusCodes[resp.StatusCode]++

	ms := float64(resp.TimeTaken.Microseconds()) / 1000.0
	r.Stats.LatencyPoints = append(r.Stats.LatencyPoints, ms)

	// Incremental avg calculation approximation or full re-calc?
	// Full recalc is expensive.
	// Running average: avg = avg + (val - avg) / count
	currentMs := float64(r.Stats.AvgLatency.Microseconds()) / 1000.0
	newAvgMs := currentMs + (ms-currentMs)/float64(r.Stats.TotalRequests)
	r.Stats.AvgLatency = time.Duration(newAvgMs * float64(time.Millisecond))

	if resp.TimeTaken < r.Stats.MinLatency {
		r.Stats.MinLatency = resp.TimeTaken
	}
	if resp.TimeTaken > r.Stats.MaxLatency {
		r.Stats.MaxLatency = resp.TimeTaken
	}
}
