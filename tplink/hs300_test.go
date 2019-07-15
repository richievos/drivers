package tplink

import (
	"github.com/reef-pi/hal"
	"testing"
)

func TestHS300Strip(t *testing.T) {
	p := NewHS300Strip("127.0.0.1:9999")
	p.cnFactory = mockConnFacctory
	d, err := HALAdapter([]byte(`{"address":"127.0.0.1:3000"}`), nil)
	if err != nil {
		t.Error(err)
	}
	if d.Metadata().Name == "" {
		t.Error("HAL metadata should not have empty name")
	}

	d1, ok := d.(hal.OutputDriver)
	if !ok {
		t.Fatal("Failed to type cast to output driver")
	}

	if len(d1.OutputPins()) != 1 {
		t.Error("Expected exactly one output pin")
	}
	pin, err := d1.OutputPin(0)
	if err != nil {
		t.Error(err)
	}
	if pin.LastState() != false {
		t.Error("Expected initial state to be false")
	}
}
