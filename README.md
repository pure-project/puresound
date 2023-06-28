# pure-sound

**A sound record/play library for go, easy to use.**

[中文](doc/README_cn.md)



### Overview

pure-sound is a simple pcm record/play library for golang, easy to use.



### Features

1. record pcm file or stream
2. play pcm file or stream
3. support Windows 7 and later (use winmm.dll)
4. support  most Linux Desktop distro (use pulseaudio)



### Usage

import:

```go
import "github.com/pure-project/puresound"
```

simple record:

```go
//do not forget check error :p

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
```

simple play:

```go
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
```

record to callback:

```go
//use the callback as a writer
var writer io.Writer = puresound.Writer(func(buf []byte) (int, error) {
	purelog.Info("recorded length: ", len(buf))
	return len(buf), nil
})

//create recorder use callback
r, err := puresound.NewRecorder(16, 16000, 1, 1600, writer)
```

play http pcm stream:

```go
//get http audio
res, err := http.Get("http://domain.live/audio.pcm")
defer res.Body.Close()

//create player
p, err := puresound.NewPlayer(16, 16000, 1, 3200, res.Body)

//start play
err = p.Start()
```



### Licence

MIT Licence

Copyright (c) 2023 pure-project team.