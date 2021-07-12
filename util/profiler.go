package util

import (
	"sync"
	"time"
)

type Profiler struct {
	curr    string
	initTime time.Time
	watches map[string]*Stopwatch
	mu sync.Mutex
}

func NewProfiler() *Profiler {
	return &Profiler{
		curr: "default",
		initTime: time.Now(),
		watches: map[string]*Stopwatch{
			"default": &Stopwatch{},
		},
	}
}

func (prof *Profiler) Start(key string) {
	prof.mu.Lock()
	defer prof.mu.Unlock()

	if key == prof.curr {
		return
	}
	prof.watches[prof.curr].Stop()
	prof.curr = key
	if _, ok := prof.watches[key]; !ok {
		prof.watches[key] = &Stopwatch{}
	}
	prof.watches[key].Start()
}

func (prof Profiler) Stop() {
	prof.mu.Lock()
	defer prof.mu.Unlock()

	prof.watches[prof.curr].Stop()
}

func (prof Profiler) ElapsedTotal() uint64 {
	return uint64(time.Since(prof.initTime).Nanoseconds())
}

func (prof Profiler) Elapsed(key string) uint64 {
	prof.mu.Lock()
	defer prof.mu.Unlock()

	if key == prof.curr {
		prof.watches[prof.curr].Stop()
	}
	sw, ok := prof.watches[key]
	if !ok {
		return 0
	}
	return uint64(sw.Elapsed().Nanoseconds())
}

type Stopwatch struct {
	startTime time.Time
	elapsed   time.Duration
	running   bool
}

func (sw *Stopwatch) Start() {
	sw.running = true
	sw.startTime = time.Now()
}

func (sw *Stopwatch) Stop() {
	elapsed := time.Since(sw.startTime)
	if !sw.running {
		return
	}
	sw.elapsed += elapsed
	sw.running = false
}

func (sw *Stopwatch) Elapsed() time.Duration {
	return sw.elapsed
}
