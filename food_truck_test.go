package stockholmfoodtrucks

import "testing"

func TestNameToHex(t *testing.T) {
	for i, tt := range []struct {
		in   string
		want string
	}{
		{"Chilibussen", "#f2900c"},
		{"El Taco Truck", "#f38ab3"},
		{"Foo Bar", "#000000"},
	} {
		if got := nameToHex(tt.in); got != tt.want {
			t.Fatalf(`[%d] nameToHex(%q) = %q, want %q`, i, tt.in, got, tt.want)
		}
	}
}
