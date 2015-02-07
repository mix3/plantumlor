package plantuml

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestPlantUML(t *testing.T) {
	p, err := NewPlantUML("")
	if err != nil {
		t.Error(err)
	}

	{
		b, err := p.Transfer("Bob->Alice : hello", PNG)
		if err != nil {
			t.Error(err)
		}
		compare(t, b, "../test_data/sample.png")
	}
	{
		b, err := p.Transfer("Bob->Alice : hello", SVG)
		if err != nil {
			t.Error(err)
		}
		compare(t, b, "../test_data/sample.svg")
	}
}

func compare(t *testing.T, got []byte, filename string) {
	expect, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(got, expect) {
		t.Error("not equal")
	}
}
