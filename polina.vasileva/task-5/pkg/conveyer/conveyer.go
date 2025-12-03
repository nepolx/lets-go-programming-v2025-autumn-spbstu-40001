package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

const Undefined = "undefined"

var ErrChannelNotFound = errors.New("chan not found")


type Conveyer struct {
	lock     sync.Mutex
	pipes    map[string]chan string
	workers  []func(context.Context) error
	capacity int
}

func New(size int) *Conveyer {
	return &Conveyer{
		pipes:    make(map[string]chan string),
		workers:  []func(context.Context) error{},
		capacity: size,
	}
}

func (c *Conveyer) ensure(name string) chan string {
	c.lock.Lock()
	defer c.lock.Unlock()

	exist := c.pipes[name]
	if exist != nil {
		return exist
	}

	created := make(chan string, c.capacity)
	c.pipes[name] = created

	return created
}

func (c *Conveyer) lookup(name string) (chan string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	ch, ok := c.pipes[name]

	return ch, ok
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input, output string,
) {
	inp := c.ensure(input)
	out := c.ensure(output)

	job := func(ctx context.Context) error {
		defer func() {
    	recover()
    	close(out)
	}()

		return fn(ctx, inp, out)
	}

	c.lock.Lock()
	c.workers = append(c.workers, job)
	c.lock.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	var inList []chan string
	for _, n := range inputs {
		inList = append(inList, c.ensure(n))
	}

	out := c.ensure(output)

	job := func(ctx context.Context) error {
		defer func() {
    		recover()
    		close(out)
		}()

		return fn(ctx, inList, out)
	}

	c.lock.Lock()
	c.workers = append(c.workers, job)
	c.lock.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inp := c.ensure(input)

	var outs []chan string
	for _, n := range outputs {
		outs = append(outs, c.ensure(n))
	}

	job := func(ctx context.Context) error {
		defer func() {
    		for _, ch := range outs {
        		func(ch chan string) {
            		defer recover()
            		close(ch)
        		}(ch)
    		}
		}()

		return fn(ctx, inp, outs)
	}

	c.lock.Lock()
	c.workers = append(c.workers, job)
	c.lock.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.lock.Lock()
	snapshot := make([]func(context.Context) error, len(c.workers))
	for i := 0; i < len(c.workers); i++ {
		snapshot[i] = c.workers[i]
	}
	c.lock.Unlock()

	group, gctx := errgroup.WithContext(ctx)

	for i := 0; i < len(snapshot); i++ {
		pos := i

		group.Go(func() error {
			return snapshot[pos](gctx)
		})
	}

	return group.Wait()
}


func (c *Conveyer) Send(input, data string) error {
	ch, ok := c.lookup(input)
	if !ok {
		return ErrChannelNotFound
	}

	defer func() {
    _ = recover()
	}()

	ch <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, ok := c.lookup(output)
	if !ok {
		return "", ErrChannelNotFound
	}

	val, ok := <-ch
	if !ok {
		return Undefined, nil
	}

	return val, nil
}
