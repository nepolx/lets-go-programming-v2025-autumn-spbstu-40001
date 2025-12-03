package conveyer

import (
	"context"
	"errors"
	"sync"
)

const Undefined = "undefined"

var (
	ErrChannelNotFound = errors.New("chan not found")
	ErrChannelFull     = errors.New("channel is full")
)

type Task func(ctx context.Context) error

type Conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	tasks    []Task
	chanSize int
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		tasks:    []Task{},
		chanSize: size,
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.chanSize)
	c.channels[name] = ch

	return ch
}

func (c *Conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.channels[name]

	return ch, ok
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input, output chan string) error,
	input, output string,
) {
	inCh := c.getOrCreateChannel(input)
	outCh := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case data, ok := <-inCh:
				if !ok {
					return nil
				}
				tmpCh := make(chan string, 1)
				tmpCh <- data
				close(tmpCh)
				if err := fn(ctx, tmpCh, outCh); err != nil {
					return err
				}
			}
		}
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inChs := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChs[i] = c.getOrCreateChannel(name)
	}

	outCh := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inCh := c.getOrCreateChannel(input)

	outChs := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChs[i] = c.getOrCreateChannel(name)
	}

	task := func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.RLock()
	tasks := append([]Task(nil), c.tasks...)
	c.mu.RUnlock()

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, t := range tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			if err := task(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(t)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		c.closeAll()
		return ctx.Err()
	case err := <-errCh:
		c.closeAll()
		return err
	case <-done:
		c.closeAll()
		return nil
	}
}

func (c *Conveyer) closeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}
}

func (c *Conveyer) Send(input string, data string) error {
	ch, ok := c.getChannel(input)
	if !ok {
		return ErrChannelNotFound
	}

	select {
	case ch <- data:
		return nil
	default:
		return ErrChannelFull
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, ok := c.getChannel(output)
	if !ok {
		return "", ErrChannelNotFound
	}

	select {
	case val, ok := <-ch:
		if !ok {
			return Undefined, nil
		}
		return val, nil
	default:
		return Undefined, nil
	}
}
