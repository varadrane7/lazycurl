package load

import (
	"fmt"
	"sort"
	"time"
)

// Stats holds the aggregated metrics of a load test.
type Stats struct {
	TotalRequests int
	ElapsedTime   time.Duration
	StatusCodes   map[int]int
	LatencyPoints []float64 // In milliseconds, for plotting
	AvgLatency    time.Duration
	MinLatency    time.Duration
	MaxLatency    time.Duration

	// Internal tracking for next window
	windowReqs    int
	windowLatency time.Duration
}

// NewStats creates a fresh stats object.
func NewStats() *Stats {
	return &Stats{
		StatusCodes: make(map[int]int),
		MinLatency:  time.Hour, // large init
	}
}

// Summary returns a formatted summary string.
func (s *Stats) Summary() string {
	return fmt.Sprintf("Total: %d | Avg: %s | Max: %s", s.TotalRequests, s.AvgLatency, s.MaxLatency)
}

// P95 calculates the 95th percentile latency (approximate if using points).
// For MVP we can just sort the huge list if it fits in memory, or keep a sample.
// We'll trust memory is fine for MVP (<10k requests).
func (s *Stats) P95() time.Duration {
	if len(s.LatencyPoints) == 0 {
		return 0
	}
	// Copy to sort so we don't mess up plotting order if we need it?
	// Actually points are chronological for plotting...
	// We should probably keep a separate list or just sort a copy.
	sorted := make([]float64, len(s.LatencyPoints))
	copy(sorted, s.LatencyPoints)
	sort.Float64s(sorted)

	idx := int(float64(len(sorted)) * 0.95)
	ms := sorted[idx]
	return time.Duration(ms * float64(time.Millisecond))
}
