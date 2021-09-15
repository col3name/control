package observe

import (
	"sync"
)

type Value struct {
	Method string
	Params []byte
}

type Observer interface {
	ID() string       // unique observer's id, attaching and detaching by this id
	Event() string    // on what event it should notified
	Notify(val Value) // notification callback
}

type Observable struct {
	mx        sync.Mutex
	observers []Observer
}

func New() *Observable {
	return &Observable{
		mx:        sync.Mutex{},
		observers: make([]Observer, 0),
	}
}

// if event is empty then event broadcasting to all observers
// if Observer.Event == '*' then this Observer handles any events
func (o *Observable) Notify(event string, val Value) {
	o.mx.Lock()
	defer o.mx.Unlock()
	for _, e := range o.observers {
		if (e.Event() == "*" || event == "" || e.Event() == event) && e.Notify != nil {
			e.Notify(val)
		}
	}
}

func (o *Observable) Register(val Observer) {
	o.mx.Lock()
	defer o.mx.Unlock()
	o.observers = append(o.observers, val)
}

func (o *Observable) Unregister(val Observer) {
	o.mx.Lock()
	defer o.mx.Unlock()
	for n, e := range o.observers {
		if e.ID() == val.ID() {
			tail := len(o.observers) - 1
			o.observers[n] = o.observers[tail]
			o.observers = o.observers[:tail]
			return
		}
	}
}
