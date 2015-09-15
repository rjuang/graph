package graph

import "time"

/*
A Port describes an instance of either an input or output port of a
graph node. Input ports are where graph nodes receive data from.
Output ports are where graph nodes send computed results through.
Each input port can only be binded to a single output port.
Each output port can be binded to multiple input ports.
*/
type Port struct {
	Data     chan interface{}
	Name     string
	bindings map[*Port]bool
	quit     chan bool
}

/*
Send the specified data through the port. This is a non-blocking call.
If the port is an input port, then the specified data is queued in the data
channel for the node function reading this input port. If the port is an output
port, then the specified data is queued multiple times up to the number of other
ports this output port is conncted to. Returns this port for chaining.
*/
func (p *Port) SendAsync(data interface{}) *Port {
	go func() {
		select {
		case p.Data <- data:
		case <-p.quit:
		}
	}()
	return p
}

/*
Binds this port (assumed to be an output port) to send data to the specified
destination port. Binding an input port to another port will give unintended
results. Returns this port for chaining.
*/
func (p *Port) BindTo(dst *Port) *Port {
	p.bindings[dst] = true
	return p
}

/*
Retrive the value sent to this port (assumed to be input port).
Reading from an output port would dequeue the results. Note this is a blocking
call.
*/
func (p *Port) Get() interface{} {
	return <-p.Data
}

/*
Start go routines associated with this port. Returns this port for chaining.
*/
func (p *Port) Start() *Port {
	// Only spawn broadcaster if output port.
	if len(p.bindings) == 0 {
		return p
	}

	go func() {
		for {
			select {
			case d := <-p.Data:

				// If data comes in, re-broadcast it to the binding
				// ports.
				for dst := range p.bindings {
					go func(dst *Port) {
						select {
						case dst.Data <- d:
						case <-p.quit:
						}
					}(dst)
				}
			case <-p.quit:
				return
			}
		}
	}()
	return p
}

/*
Shutdown this port. This is a blocking call. Returns this port for chaining.
*/
func (p *Port) Shutdown() *Port {
	for {
		select {
		case p.quit <- true:
			// Keep looping until quit exhausted determined by timeout.
		case <-time.After(100 * time.Millisecond):
			return p
		}
	}
}

/*
Construct and return a new port.
*/
func CreatePort(name string) *Port {
	p := &Port{
		Data:     make(chan interface{}),
		Name:     name,
		bindings: make(map[*Port]bool),
		quit:     make(chan bool),
	}
	return p
}
