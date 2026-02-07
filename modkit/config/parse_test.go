package config_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/go-modkit/modkit/modkit/config"
)

func TestParseHelpers(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		got, err := config.ParseString("hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "hello" {
			t.Fatalf("got %q", got)
		}
	})

	t.Run("int", func(t *testing.T) {
		got, err := config.ParseInt("42")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 42 {
			t.Fatalf("got %d", got)
		}

		if _, err := config.ParseInt("nope"); err == nil {
			t.Fatalf("expected parse error")
		}
	})

	t.Run("float64", func(t *testing.T) {
		got, err := config.ParseFloat64("3.5")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 3.5 {
			t.Fatalf("got %f", got)
		}

		if _, err := config.ParseFloat64("nope"); err == nil {
			t.Fatalf("expected parse error")
		}
	})

	t.Run("bool", func(t *testing.T) {
		got, err := config.ParseBool("true")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !got {
			t.Fatalf("got false")
		}

		if _, err := config.ParseBool("not-bool"); err == nil {
			t.Fatalf("expected parse error")
		}
	})

	t.Run("duration", func(t *testing.T) {
		got, err := config.ParseDuration("90s")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 90*time.Second {
			t.Fatalf("got %v", got)
		}

		if _, err := config.ParseDuration("1 hour"); err == nil {
			t.Fatalf("expected parse error")
		}
	})

	t.Run("csv", func(t *testing.T) {
		got, err := config.ParseCSV(" a, ,b,c ,, ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []string{"a", "b", "c"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v want %v", got, want)
		}
	})
}
