package nanux

import (
	"errors"

	"github.com/nanux-io/nanux/handler"
	"github.com/nanux-io/nanux/transporter"
)

// Nanux is the base struct to use nanux. It contains the used listeners
// (nats, http, ...) and a context which can hold for example the instance
// of the DB to used.
type Nanux struct {
	// L is the instance of listener used by nanux
	L transporter.Listener
	// Ctx is an a nanux scoped context
	Ctx           interface{}
	isActionAdded bool
	errorHandler  handler.ManageError
}

// Handle defines the action to execute when the given route is reached on the
// listener.
// A route is an http route in the case of an http listener, its a channel
// subscription in the case of a nats listener etc...
func (n *Nanux) Handle(route string, a handler.Action) error {
	n.isActionAdded = true
	fn := func(req handler.Request) ([]byte, error) {
		return a.Fn(&n.Ctx, req)
	}

	listenerAction := handler.ListenerAction{
		Fn:   fn,
		Opts: a.Opts,
	}
	err := n.L.HandleAction(route, listenerAction)

	return err
}

// Listen start to listen whith the listener (nats, http, etc...) attached to
// Nanux.
func (n *Nanux) Listen() error {
	return n.L.Listen()
}

// Close the listener connection
func (n *Nanux) Close() error {
	return n.L.Close()
}

// HandleError specify error handler which must be called when an action return
// en error
func (n *Nanux) HandleError(errHandler handler.ManageError) error {
	if n.isActionAdded == true {
		return errors.New("Can not add error handler after adding action to handle")
	}

	return n.L.HandleError(errHandler)
}

// New create a new Nanux for nanux
func New(listener transporter.Listener, ctx interface{}) *Nanux {
	return &Nanux{
		L:   listener,
		Ctx: ctx,
	}
}
