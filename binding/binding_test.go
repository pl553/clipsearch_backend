package binding

import (
	"testing"
)

func TestShouldBind(t *testing.T) {
	t.Run("decode error", func(t *testing.T) {
		src := map[string][]string{
			"field": {"value"},
		}
		dst := struct {
			Field int `schema:"field"`
		}{}
		err := ShouldBind(&dst, src)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("int field validation error", func(t *testing.T) {
		src := map[string][]string{
			"fieldint": {"-1"},
		}
		dst := struct {
			Field int `schema:"fieldint" validate:"min=0"`
		}{}
		err := ShouldBind(&dst, src)
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("no errors", func(t *testing.T) {
		src := map[string][]string{
			"a": {"10"},
			"b": {"-102"},
		}
		dst := struct {
			A int `schema:"a" validate:"min=10,max=10"`
			B int `schema:"b" validate:"max=-100"`
		}{}
		err := ShouldBind(&dst, src)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if dst.A != 10 || dst.B != -102 {
			t.Errorf("Expected dst.A = %v, dst.B = %v, got %v %v", 10, -102, dst.A, dst.B)
		}
	})
}
