package graph

/*
Graph structure to keep track of nodes associated with different graphs. This
allows developers to start/shutdown nodes of specific graphs.
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
Start the graph.
*/
func (g *Graph) Start() {
	for _, n := range g.nodes {
		n.Start()
	}
}

/*
Shutdown the graph.
*/
func (g *Graph) Shutdown() {
	for _, n := range g.nodes {
		n.Shutdown()
	}
}

/*
Construct an empty graph.
*/
func CreateGraph() *Graph {
	return &Graph{
		nodes: []*Node{},
	}
}
