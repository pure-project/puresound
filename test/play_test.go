package test

import (
	"puresound"
	"testing"
)

func TestDefaultPlay(t *testing.T) {
	dev, err := puresound.DefaultPlay()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("name=%s handle=%+v", dev.Name(), dev.Handle())
}

func TestListPlay(t *testing.T) {
	devs, err := puresound.ListPlay()
	if err != nil {
		t.Fatal(err)
	}

	for _, dev := range devs {
		t.Logf("name=%s handle=%+v", dev.Name(), dev.Handle())
	}
}