package conveyer

import (
	"context"
	"errors"
	"fmt"
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
		lock:     sync.Mutex{},
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
	functionn func(ctx context.Context, input chan string, output chan string) error,
	input, output string,
) {
	inp := c.ensure(input)
	out := c.ensure(output)

	job := func(ctx context.Context) error {
		defer close(out)

		return functionn(ctx, inp, out)
	}

	c.lock.Lock()
	c.workers = append(c.workers, job)
	c.lock.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	functionn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inList := make([]chan string, 0, len(inputs))
	for _, n := range inputs {
		inList = append(inList, c.ensure(n))
	}

	out := c.ensure(output)

	job := func(ctx context.Context) error {
		defer close(out)

		return functionn(ctx, inList, out)
	}

	c.lock.Lock()
	c.workers = append(c.workers, job)
	c.lock.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	functionn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inp := c.ensure(input)

	outs := make([]chan string, 0, len(outputs))
	for _, n := range outputs {
		outs = append(outs, c.ensure(n))
	}

	job := func(ctx context.Context) error {
		defer func() {
			for _, ch := range outs {
				defer close(ch)
			}
		}()

		return functionn(ctx, inp, outs)
	}

	c.lock.Lock()
	c.workers = append(c.workers, job)
	c.lock.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.lock.Lock()
	snapshot := make([]func(context.Context) error, len(c.workers))
	copy(snapshot, c.workers)
	c.lock.Unlock()

	group, gctx := errgroup.WithContext(ctx)

	for i := range snapshot {
		pos := i

		group.Go(func() error {
			return snapshot[pos](gctx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input, data string) error {
	channel, exists := c.lookup(input)
	if !exists {
		return ErrChannelNotFound
	}

	defer func() {
		_ = recover()
	}()

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	channel, exists := c.lookup(output)
	if !exists {
		return "", ErrChannelNotFound
	}

	val, ok := <-channel
	if !ok {
		return Undefined, nil
	}

	return val, nil
}
