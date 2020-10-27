package parallel

import "sync"

// Parallel provides parallel computation functions by managing goroutines
// explicitly.
type Parallel struct {
	daemons   chan func()
	waitGroup *sync.WaitGroup
}

// NewParallel creates a parallel object.
func NewParallel(maxGoroutines int) *Parallel {
	return &Parallel{make(chan func(), maxGoroutines), &sync.WaitGroup{}}
}

func (p *Parallel) addDaemon(f func()) {
	p.waitGroup.Add(1)

	p.daemons <- func() {
		f()
		p.waitGroup.Done()
	}
}

// Run runs daemons util all of them finish.
func (p *Parallel) Run() {
	go func() {
		for f := range p.daemons {
			go f()
		}
	}()

	p.waitGroup.Wait()
}

// UnorderedMap maps.
func (p *Parallel) UnorderedMap(fs chan func() interface{}) chan interface{} {
	xs := make(chan interface{}, cap(fs))

	for f := range fs {
		p.addDaemon(func() {
			xs <- f()
		})
	}

	return xs
}
