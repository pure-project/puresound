package test

import (
	"github.com/pure-project/puresound"
	"os"
	"testing"
	"time"
)

func TestRecord(t *testing.T) {
	f, err := os.OpenFile("test.pcm", os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0666)
	if err != nil {
		t.Fatal(err)
	}

	r, err := puresound.NewRecorder(16, 16000, 1, 1600, f)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Start()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	r.Stop()


}