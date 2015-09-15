package graph

import "testing"

func TestCreatePort(t *testing.T) {
	p := CreatePort("myport")
	if p.Name != "myport" {
		t.Errorf("Expecting name to be '%s'. Got '%s'.", "myport", p.Name)
	}

	go func() {
		p.SendAsync("test")
	}()

	v, ok := p.Get().(string)
	if !ok {
		t.Errorf("Incorrect type retrieved in Get.")
	}
	if v != "test" {
		t.Errorf("Expected '%s'. Got '%s'.", "test", v)
	}
	// Start and shutdown should do nothing since no bindings.
	p.Start()
	p.Shutdown()
}

func TestSingleBinding(t *testing.T) {
	p1 := CreatePort("myport1")
	p2 := CreatePort("myport2")

	p1.BindTo(p2)
	p1.Start()
	p2.Start()

	p1.SendAsync("hello world")
	v2, ok2 := p2.Get().(string)

	if !ok2 {
		t.Errorf("Incorrect type retrieved in Get calls.")
	}

	if v2 != "hello world" {
		t.Errorf("v2: %s. Expected '%s'.", v2, "hello world")
	}

	p1.Shutdown()
	p2.Shutdown()
}

func TestMultiBinding(t *testing.T) {
	p1 := CreatePort("myport1")
	p2 := CreatePort("myport2")
	p3 := CreatePort("myport3")

	p1.BindTo(p2)
	p1.BindTo(p3)
	p1.Start()
	p2.Start()
	p3.Start()

	p1.SendAsync("hello world")
	v2, ok2 := p2.Get().(string)
	v3, ok3 := p3.Get().(string)

	if !ok2 || !ok3 {
		t.Errorf("Incorrect type retrieved in Get calls.")
	}

	if v2 != "hello world" || v3 != "hello world" {
		t.Errorf("v2: %s. v3: %s. Expected '%s'.", v2, v3, "hello world")
	}

	p1.Shutdown()
	p2.Shutdown()
	p3.Shutdown()
}
