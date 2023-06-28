package test

import (
	"os"
	"puresound"
	"testing"
	"time"
)

func TestPlay(t *testing.T) {
	f, err := os.Open("test.pcm")
	if err != nil {
		t.Fatal(err)
	}

	p, err := puresound.NewPlayer(16, 16000, 1, 3200, f)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Start()
	if err != nil {
		t.Fatal(err)
	}

	//time.Sleep(5 * time.Second)

	for p.Playing() {
		time.Sleep(100 * time.Millisecond)
	}

	err = p.Stop()
	if err != nil {
		t.Fatal(err)
	}

	err = p.Close()
	if err != nil {
		t.Fatal(err)
	}
}