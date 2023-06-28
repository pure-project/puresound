package test

import (
	"puresound"
	"testing"
)

func TestDefaultRecord(t *testing.T) {
	dev, err := puresound.DefaultRecord()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("name=%s handle=%v", dev.Name(), dev.Handle())
}

func TestListRecord(t *testing.T) {
	devs, err := puresound.ListRecord()
	if err != nil {
		t.Fatal(err)
	}

	for _, dev := range devs {
		t.Logf("name=%s handle=%v", dev.Name(), dev.Handle())
	}
}