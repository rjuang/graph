package graph

import "testing"

func assertIntEq(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected: %d != Actual: %d", expected, actual)
	}
}

func TestCreateNode(t *testing.T) {
	fn := func(in DataMap, params DataMap) DataMap {
		i1 := in["in1"].(int)
		i2 := in["in2"].(int)

		return DataMap{
			"out1": i1 + i2,
			"out2": i1 - i2,
		}
	}

	node := CreateNode([]string{"in1", "in2"}, []string{"out1", "out2"}, fn)
	node.Start()
	node.In["in1"].SendAsync(10)
	node.In["in2"].SendAsync(8)

	r1 := node.Out["out1"].Get().(int)
	r2 := node.Out["out2"].Get().(int)

	node.Shutdown()
	assertIntEq(t, 18, r1)
	assertIntEq(t, 2, r2)
}

func TestMultipleNodes(t *testing.T) {

	fn1 := func(in DataMap, params DataMap) DataMap {
		i1 := in["in1"].(int)
		i2 := in["in2"].(int)

		return DataMap{
			"out1": i1 + i2,
			"out2": i1 - i2,
		}
	}

	fn2 := func(in DataMap, params DataMap) DataMap {
		i1 := in["in1"].(int)
		i2 := in["in2"].(int)

		return DataMap{
			"out": i1 - i2,
		}
	}

	node1 := CreateNode([]string{"in1", "in2"}, []string{"out1", "out2"}, fn1)
	node2 := CreateNode([]string{"in1", "in2"}, []string{"out"}, fn2)
	node3 := CreateNode([]string{"in1", "in2"}, []string{"out"}, fn2)
	node1.Out["out1"].BindTo(node2.In["in1"])
	node1.Out["out2"].BindTo(node2.In["in2"])
	node1.Out["out1"].BindTo(node3.In["in2"])
	node1.Out["out2"].BindTo(node3.In["in1"])

	node1.Start()
	node2.Start()
	node3.Start()
	node1.In["in1"].SendAsync(10)
	node1.In["in2"].SendAsync(8)

	r2 := node2.Out["out"].Get().(int)
	r3 := node3.Out["out"].Get().(int)

	node1.Shutdown()
	node2.Shutdown()
	node3.Shutdown()

	assertIntEq(t, 16, r2)
	assertIntEq(t, -16, r3)
}

func TestMultipleNodesSetInputs(t *testing.T) {

	fn1 := func(in DataMap, params DataMap) DataMap {
		i1 := in["in1"].(int)
		i2 := in["in2"].(int)

		return DataMap{
			"out1": i1 + i2,
			"out2": i1 - i2,
		}
	}

	fn2 := func(in DataMap, params DataMap) DataMap {
		i1 := in["in1"].(int)
		i2 := in["in2"].(int)

		return DataMap{
			"out": i1 - i2,
		}
	}

	node1 := CreateNode([]string{"in1", "in2"}, []string{"out1", "out2"}, fn1)
	node2 := CreateNode([]string{"in1", "in2"}, []string{"out"}, fn2)
	node3 := CreateNode([]string{"in1", "in2"}, []string{"out"}, fn2)

	node2.SetInputs(map[string]*Port{
		"in1": node1.Out["out1"],
		"in2": node1.Out["out2"],
	})

	node3.SetInputs(map[string]*Port{
		"in1": node1.Out["out2"],
		"in2": node1.Out["out1"],
	})

	node1.Start()
	node2.Start()
	node3.Start()
	node1.In["in1"].SendAsync(10)
	node1.In["in2"].SendAsync(8)

	r2 := node2.Out["out"].Get().(int)
	r3 := node3.Out["out"].Get().(int)

	node1.Shutdown()
	node2.Shutdown()
	node3.Shutdown()

	assertIntEq(t, 16, r2)
	assertIntEq(t, -16, r3)
}
