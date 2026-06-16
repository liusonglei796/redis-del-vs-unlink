package main

import (
	"testing"
)

func TestFormatLog(t *testing.T) {
	expected := "Duration: 150 us | Cmd: DEL"
	result := FormatLog(150, "DEL")
	if result != expected {
		t.Errorf("FormatLog(150, \"DEL\") = %q; want %q", result, expected)
	}
}
