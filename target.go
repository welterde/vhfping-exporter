package main

import (
	"log"
	"sync"
	
	"github.com/VictoriaMetrics/metrics"
)

// Target is the internal state about a target
type Target struct {
	sync.Mutex
	spec         TargetSpec
	vset         *metrics.Set
//	measurements *Measurements
}

func NewTarget(ts TargetSpec) *Target {
	t := Target{
		spec:     ts,
		vset: metrics.NewSet(),
	}

	log.Println("new target:", ts.host)

	return &t
}

// fping_sent_count
// fping_lost_count
// fping_rtt_sum
// fping_rtt_count
// fping_rtt_hist

func (t *Target) AddMeasurements(m Measurements) {
	t.Lock()
	//t.measurements = &m
	sent_count := t.vset.GetOrCreateGauge("vhfping_sent_count", nil)
	sent_count.Set(float64(m.GetSentCount()))
	lost_count := t.vset.GetOrCreateGauge("vhfping_lost_count", nil)
	lost_count.Set(float64(m.GetLostCount()))
	median_rtt := t.vset.GetOrCreateGauge("vhfping_median_rtt", nil)
	median_rtt.Set(float64(m.GetMedian()))
	rtt_hist := t.vset.GetOrCreateHistogram("vhfping_rtt")
	//rtt_hist.Reset()
	for i := range m.rtt {
		if !m.lost[i] {
			rtt_hist.Update(m.rtt[i])
		}
	}
	t.Unlock()
}
