package graph

import "time"

type DataMap map[string]interface{}

/*
Function specified by user to define the computational work of the node.
The node function is called each time a complete set of node inputs is received.

Args:
  in contains a mapping of the port name to the actual data.
  params contains a mapping of the static parameter name and value.

Returns:
  a mapping where the key represents the output port name and value contains the
  output data to send to the respective port.
*/
type NodeFunc func(in DataMap, params DataMap) DataMap

/*
Node represents a node in a computational graph. Each node of the graph can do
work on inputs provided to it and output the resulting values to other nodes
connected to it via the output ports. Each node has a set of parameters
that represents constant settings (e.g. threshold value, model parameters, etc).
*/
type Node struct {
	In     map[string]*Port
	Out    map[string]*Port
	Params DataMap
	fn     NodeFunc
	quit   chan bool
	stats  map[string]float64
}

/*
Start this node and kickoff all go routines to handle computation for this node.
Returns this node for chaining.
*/
func (n *Node) Start() *Node {
	for _, p := range n.Out {
		p.Start()
	}

	go func() {
		for {
			// Needed in case no inputs to this node.
			select {
			case <-n.quit:
				return
			default:
			}

			inputs := make(DataMap)
			for k, p := range n.In {
				select {
				case inputs[k] = <-p.Data:
				case <-n.quit:
					return
				}
			}

			// Input complete. Trigger function.
			n.stats["count"] += 1
			startTime := time.Now()
			outputs := n.fn(inputs, n.Params)
			elapsedMs := time.Since(startTime).Seconds() * 1000.0
			n.stats["last_runtime"] = elapsedMs

			w := 1.0 / n.stats["count"]
			n.stats["avg_runtime"] = ((1.0 - w) * n.stats["avg_runtime"]) + (w * elapsedMs)

			for k, v := range outputs {
				if p, ok := n.Out[k]; ok {
					p.SendAsync(v)
				}
			}
		}
	}()

	return n
}

/*
Shutdown the node. This is a blocking call. Returns this node for chaining.
*/
func (n *Node) Shutdown() *Node {
	for {
		quit := false
		select {
		case n.quit <- true:
		case <-time.After(100 * time.Millisecond):
			quit = true
		}
		if quit {
			break
		}
	}
	for _, p := range n.Out {
		p.Shutdown()
	}
	return n
}

/*
Connect input ports to the specified ports. Input is a map containing the input
ports and the corresponding output ports to connect to. Returns this node for
chaining.
*/
func (n *Node) SetInputs(bindings map[string]*Port) *Node {
	for i, o := range bindings {
		o.BindTo(n.In[i])
	}
	return n
}

/*
Construct a node.

Example usage:

    fn := func(in DataMap, params DataMap) DataMap {
	// These would need to fetch synchronously since function computation depends on these values
	bytes1 := in["img1"].([]byte)
	bytes2 := in["img2"].([]byte)

	result1 := ... // some computation here
	result2 := ... // some computation here

	return DataMap {
	    "diff": result1,
	    "mask": result2,
	}
    }

    graph.CreateNode([]string{"img1", "img2"}, []string{"diff", "mask"}, fn)
*/
func CreateNode(in []string, out []string, fn NodeFunc) *Node {
	inputPorts := make(map[string]*Port)
	for _, name := range in {
		inputPorts[name] = CreatePort(name)
	}

	outputPorts := make(map[string]*Port)
	for _, name := range out {
		outputPorts[name] = CreatePort(name)
	}

	return &Node{
		In:     inputPorts,
		Out:    outputPorts,
		Params: make(DataMap),
		fn:     fn,
		quit:   make(chan bool),
		stats:  make(map[string]float64),
	}
}
