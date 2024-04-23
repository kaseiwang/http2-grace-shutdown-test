### steps to reproduce

```sh
go run .
```

### my test result

before patch:
```sh
tcp connected
sleep for 15 seconds, exec "nginx -s reload" now
sleep done, sending http2 preface
http2 init done, sending requests. 10 RPS.
recv frame: [FrameHeader WINDOW_UPDATE len=4]
recv frame: [FrameHeader SETTINGS flags=ACK len=0]
recv frame: [FrameHeader SETTINGS flags=ACK len=0]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=1 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=3 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=5 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=7 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=9 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=11 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=13 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=15 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=17 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=19 len=57]

.......

recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=1991 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=1993 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=1995 len=57]
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=1997 len=57]
recv frame: [FrameHeader GOAWAY len=8]
server says go away
recv frame: [FrameHeader HEADERS flags=END_STREAM|END_HEADERS stream=1999 len=57]
read err: EOF
```


after patch:
```sh
tcp connected
sleep for 15 seconds, exec "nginx -s reload" now
sleep done, sending http2 preface
http2 init done, sending requests. 10 RPS.
recv frame: [FrameHeader WINDOW_UPDATE len=4]
recv frame: [FrameHeader SETTINGS flags=ACK len=0]
recv frame: [FrameHeader GOAWAY len=8]
server says go away
read err: EOF
```