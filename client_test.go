package stockholmfoodtrucks

import (
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	in, out := "baz", "bar"

	os.Setenv("ENVSTR", out)

	if got := Env("ENVSTR", in); got != out {
		t.Errorf(`String("ENVSTR", "%v") = %v, want %v`, in, got, out)
	}
}

func TestEnvDefault(t *testing.T) {
	in, out := "baz", "baz"

	if got := Env("ENVSTR_DEFAULT", in); got != out {
		t.Errorf(`String("ENVSTR_DEFAULT", "%v") = %v, want %v`, in, got, out)
	}
}
