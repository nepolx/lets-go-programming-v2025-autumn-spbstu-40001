package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrNoDecorator = errors.New("can't be decorated")
	ErrNoOutputs   = errors.New("no output channels provided for separator")
)

func PrefixDecoratorFunc(ctx context.Context, inChan, outChan chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-inChan:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outChan <- val:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, inChan chan string, outChans []chan string) error {
	if len(outChans) == 0 {
		return ErrNoOutputs
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-inChan:
			if !ok {
				return nil
			}

			target := outChans[index%len(outChans)]
			index++

			select {
			case <-ctx.Done():
				return ctx.Err()
			case target <- val:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inChans []chan string, outChan chan string) error {
	var wg sync.WaitGroup

	for _, ch := range inChans {
		wg.Add(1)
		go func(c chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-c:
					if !ok {
						return
					}
					if strings.Contains(val, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case outChan <- val:
					}
				}
			}
		}(ch)
	}

	wg.Wait()
	return nil
}
