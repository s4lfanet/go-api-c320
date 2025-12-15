package app

import (
	"testing"
)

func TestNew(t *testing.T) {
	app := New()

	if app == nil {
		t.Error("Expected non-nil app instance")
	}

	if app.router != nil {
		t.Error("Expected router to be nil before Start is called")
	}
}

func TestNew_ReturnsAppStruct(t *testing.T) {
	app1 := New()
	app2 := New()

	if app1 == nil || app2 == nil {
		t.Error("Expected both app instances to be non-nil")
	}

	// Each call should return a new instance
	if app1 == app2 {
		t.Error("Expected different app instances")
	}
}
