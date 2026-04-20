package internal

import (
	"strings"
	"testing"
)

func TestDim(t *testing.T) {
	input := "low visibility"
	result := Dim(input)

	// Check if it contains the input text
	if !strings.Contains(result, input) {
		t.Errorf("Dim() result does not contain input text")
	}

	// Check if it starts with the Dim color and ends with Reset
	if !strings.HasPrefix(result, ColorDim) {
		t.Errorf("Dim() result should start with ColorDim")
	}
	if !strings.HasSuffix(result, ColorReset) {
		t.Errorf("Dim() result should end with ColorReset")
	}
}

func TestIndigo(t *testing.T) {
	input := "bright indigo"
	result := Indigo(input)

	if !strings.Contains(result, input) {
		t.Errorf("Indigo() result does not contain input text")
	}

	if !strings.HasPrefix(result, ColorIndigo) {
		t.Errorf("Indigo() result should start with ColorIndigo")
	}
}

func TestStyledBar(t *testing.T) {
	input := "MY TASKS"
	result := StyledBar(input)

	// Check if input is present
	if !strings.Contains(result, input) {
		t.Errorf("StyledBar() result does not contain input text")
	}

	// Check for the background color and the newline at the end
	if !strings.Contains(result, BgIndigo) {
		t.Errorf("StyledBar() missing background color")
	}
	if !strings.HasSuffix(result, "\n") {
		t.Errorf("StyledBar() must end with a newline")
	}
}

func TestFancyBar(t *testing.T) {
	title := "TODO APP"
	version := "v1.0.0"
	result := FancyBar(title, version)

	if !strings.Contains(result, title) {
		t.Errorf("FancyBar() missing title")
	}
	if !strings.Contains(result, version) {
		t.Errorf("FancyBar() missing version")
	}
	if !strings.Contains(result, "⚡") {
		t.Errorf("FancyBar() missing icon")
	}
	if !strings.HasSuffix(result, "\n") {
		t.Errorf("FancyBar() must end with a newline")
	}
}

func TestMegaLogo(t *testing.T) {
	result := MegaLogo()

	// Check for a unique part of the ASCII art
	if !strings.Contains(result, "██╗") {
		t.Errorf("MegaLogo() ASCII art seems broken or missing")
	}

	// Ensure color is applied
	if !strings.HasPrefix(result, ColorIndigo) {
		t.Errorf("MegaLogo() should start with ColorIndigo")
	}

	// Ensure it resets and ends with newline
	if !strings.HasSuffix(result, ColorReset+"\n") {
		t.Errorf("MegaLogo() should reset color and end with newline")
	}
}
