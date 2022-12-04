package go_then

import (
	"errors"
	"sync"
	"time"
)

const PROMISE_TIMEOUT_IN_SECS = 60

type Promise struct {
	resolveChannel chan any
	wg             sync.WaitGroup
	errorHandler   func(err error)
	timeOutInSecs  int
}

type Config struct {
	TimeOutInSecs int
}

func New(config *Config) *Promise {

	p := new(Promise)
	p.resolveChannel = make(chan any, 1)

	if config != nil && config.TimeOutInSecs > 0 {
		p.timeOutInSecs = config.TimeOutInSecs
	}

	return p
}

func (p *Promise) Resolve(i any) {
	p.resolveChannel <- i
}

func (p *Promise) Reject(i error) {
	p.errorHandler(i)
	p.wg.Done()
}

func (p *Promise) Execute(task func()) *Promise {
	p.wg.Add(1)
	go task()

	return p
}

func (p *Promise) Wait() {
	p.wg.Wait()
}

func (p *Promise) Then(callback func(i any)) *Promise {
	go p.thenExecutor(callback)
	return p
}

func (p *Promise) thenExecutor(callback func(i any)) {
	for i := 0; i < p.timeOutInSecs; i++ {
		select {
		case o := <-p.resolveChannel:
			callback(o)
			p.wg.Done()
			return
		default:
		}
		time.Sleep(time.Second * 1)
	}

	p.errorHandler(errors.New("timeout"))
	p.wg.Done()
}

func (p *Promise) Catch(errorHandler func(err error)) {
	p.errorHandler = errorHandler
}
