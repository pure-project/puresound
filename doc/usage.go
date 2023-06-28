package main

import (
	"io"
	"net/http"
	"os"
	"puresound"
	"time"
)

func simpleRecord() {
	//create output pcm file
	out, err := os.OpenFile("output.pcm", os.O_CREATE | os.O_WRONLY, 0666)

	//create recorder
	r, err := puresound.NewRecorder(16, 16000, 1, 1600, out)

	//start record
	err = r.Start()

	//record duration
	time.Sleep(10 * time.Second)

	//stop record
	err = r.Stop()

	//close recorder
	r.Close()

	_ = err
}

func simplePlay() {
	//open input pcm file
	in, err := os.Open("input.pcm")

	//create player
	p, err := puresound.NewPlayer(16, 16000, 1, 1600, in)

	//start play
	err = p.Start()

	//pause play
	err = p.Pause()

	//resume play
	err = p.Resume()

	//wait play over
	for p.Playing() {
		time.Sleep(50 * time.Millisecond)
	}

	//stop play
	err = p.Stop()

	//close player
	p.Close()

	_ = err
}

func playStream() {
	//get http audio
	res, err := http.Get("http://domain.live/audio.pcm")

	//create player
	p, err := puresound.NewPlayer(16, 16000, 1, 3200, res.Body)

	//start play
	err = p.Start()

	_ = err
}

func recordCallback() {
	//use the callback as a writer
	var writer io.Writer = puresound.Writer(func(buf []byte) (int, error) {
		println("recorded length: ", len(buf))
		return len(buf), nil
	})

	//create recorder use callback
	r, err := puresound.NewRecorder(16, 16000, 1, 1600, writer)
}