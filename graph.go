package graph

/*
Graph structure to keep track of nodes associated with different graphs. This
allows developers to start/shutdown nodes of specific graphs.

Example usage:

g := CreateGraph()

fn1 := func(in DataMap, params DataMap) DataMap {
    in1 := in["in1"].([]byte)
    in2 := in["in2"].([]byte)

    // ... Do some computation over in1 and in2
    data := // do some computation...

    return DataMap {
	"out": data,
    }
}

g.Node([]string{"in1", "in2"}, []string{"out1", "out2"}, fn1)
*/
type Graph struct {
	nodes []*Node
}

/*
Create and return new node for this graph.
*/
func (g *Graph) Node(in []string, out []string, fn NodeFunc) *Node {
	n := CreateNode(in, out, fn)
	g.nodes = append(g.nodes, n)
	return n
}

/*
Start the graph. Returns this graph for possible chaining.
*/
func (g *Graph) Start() *Graph {
	for _, n := range g.nodes {
		n.Start()
	}
	return g
}

/*
Shutdown the graph. Returns this graph for possible chaining.
*/
func (g *Graph) Shutdown() *Graph {
	for _, n := range g.nodes {
		n.Shutdown()
	}
	return g
}

/*
Construct an empty graph.
*/
func CreateGraph() *Graph {
	return &Graph{
		nodes: []*Node{},
	}
}
