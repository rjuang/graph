package graph

import "testing"

func TestGraphWithNoElements(t *testing.T) {
	g := CreateGraph()
	g.Start()
	g.Shutdown()
}

func TestGraphWithMultipleElements(t *testing.T) {
	fn1 := func(in map[string]interface{}, params map[string]interface{}) map[string]interface{} {
		i1 := in["in1"].(int)
		i2 := in["in2"].(int)

		return map[string]interface{}{
			"out1": i1 + i2,
			"out2": i1 - i2,
		}
	}

	fn2 := func(in map[string]interface{}, params map[string]interface{}) map[string]interface{} {
		i1 := in["in1"].(int)
		i2 := in["in2"].(int)

		return map[string]interface{}{
			"out": i1 - i2,
		}
	}

	g := CreateGraph()
	node1 := g.Node([]string{"in1", "in2"}, []string{"out1", "out2"}, fn1)
	node2 := g.Node([]string{"in1", "in2"}, []string{"out"}, fn2)
	node3 := g.Node([]string{"in1", "in2"}, []string{"out"}, fn2)

	node1.Out["out1"].BindTo(node2.In["in1"])
	node1.Out["out2"].BindTo(node2.In["in2"])
	node1.Out["out1"].BindTo(node3.In["in2"])
	node1.Out["out2"].BindTo(node3.In["in1"])

	g.Start()
	node1.In["in1"].SendAsync(10)
	node1.In["in2"].SendAsync(8)
	r2 := node2.Out["out"].Get().(int)
	r3 := node3.Out["out"].Get().(int)
	g.Shutdown()

	assertIntEq(t, 16, r2)
	assertIntEq(t, -16, r3)
}
