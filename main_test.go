package main

import "testing"

func TestParse(t *testing.T) {
	cases := map[string]int{
		"5.5":   55,
		"0.0":   0,
		"-0.5":  -5,
		"-99.9": -999,
		"99.9":  999,
	}
	for in, want := range cases {
		if got := customTemperatureParser([]byte(in)); got != want {
			t.Errorf("%s: got %d, want %d", in, got, want)
		}
	}
}
