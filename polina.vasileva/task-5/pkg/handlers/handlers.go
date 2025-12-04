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
			return nil
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
			case outChan <- val:
			case <-ctx.Done():
				return nil
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
			return nil
		case val, ok := <-inChan:
			if !ok {
				return nil
			}

			target := outChans[index%len(outChans)]
			index++

			select {
			case target <- val:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inChans []chan string, outChan chan string) error {
	var waitGroup sync.WaitGroup
	for _, channel := range inChans {
		waitGroup.Add(1)

		worker := func(inputChannel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-inputChannel:
					if !ok {
						return
					}

					if strings.Contains(val, "no multiplexer") {
						continue
					}

					select {
					case outChan <- val:
					case <-ctx.Done():
						return
					}
				}
			}
		}
		go worker(channel)
	}

	waitGroup.Wait()

	return nil
}
