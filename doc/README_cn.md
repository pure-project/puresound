# pure-sound

**一个用于录制、播放声音的golang库，简单易用。**





### 前言

pure-sound是一个简单的golang pcm录制、播放库，简单易用。



### 功能

1. 录制pcm文件或数据流
2. 播放pcm文件或数据流
3. 支持Windows 7及更新版本 (使用winmm.dll)
4. 支持大多数Linux桌面发行版 (使用pulseaudio)



### 用法

导入包：

```go
import "github.com/pure-project/puresound"
```

简单录制：

```go
//记得处理错误

//创建输出pcm文件
out, err := os.OpenFile("output.pcm", os.O_CREATE | os.O_WRONLY, 0666)

//创建录音实例
r, err := puresound.NewRecorder(16, 16000, 1, 1600, out)

//开始录制
err = r.Start()

//控制录制时长
time.Sleep(10 * time.Second)

//停止录制
err = r.Stop()

//关闭录音实例
r.Close()
```

简单播放:

```go
//打开输入pcm文件
in, err := os.Open("input.pcm")

//创建播放实例
p, err := puresound.NewPlayer(16, 16000, 1, 1600, in)

//开始播放
err = p.Start()

//暂停播放
err = p.Pause()

//恢复播放
err = p.Resume()

//检测播放状态
for p.Playing() {
	time.Sleep(50 * time.Millisecond)
}

//停止播放
err = p.Stop()

//关闭播放实例
p.Close()
```

录制到回调:

```go
//使用回调函数作为写入器
var writer io.Writer = puresound.Writer(func(buf []byte) (int, error) {
	purelog.Info("recorded length: ", len(buf))
	return len(buf), nil
})

//使用回调创建录制器
r, err := puresound.NewRecorder(16, 16000, 1, 1600, writer)
```

播放HTTP PCM流:

```go
//HTTP GET获取数据流
res, err := http.Get("http://domain.live/audio.pcm")
defer res.Body.Close()

//创建播放器
p, err := puresound.NewPlayer(16, 16000, 1, 3200, res.Body)

//开始播放
err = p.Start()
```



### 许可

MIT许可证

版权所有 (c) 2023 pure-project团队。