package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

var globalCloser = New()

func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func Wait() {
	globalCloser.Wait()
}

func CloseAll() {
	globalCloser.CloseAll()
}

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []func() error
	done  chan struct{}
}

func New(sig ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(sig) > 0 {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, sig...)
		<-ch
		signal.Stop(ch)
		c.CloseAll()
	}
	return c
}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func (c *Closer) Wait() {
	<-c.done
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, fn := range funcs {
			go func(fn func() error) {
				errs <- fn()
			}(fn)
		}

		for i := 0; i < len(errs); i++ {
			if err := <-errs; err != nil {
				log.Printf("error returned from Closer: %s", err.Error())
			}
		}
	})

}
