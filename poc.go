package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

const ServerAddr = "172.17.0.2:80"

var clientPreface = []byte(http2.ClientPreface)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

type h2c struct {
	conn                 net.Conn
	framer               *http2.Framer
	MaxConcurrentStreams uint32
}

func (c *h2c) sendPreface() bool {
	n, err := c.conn.Write(clientPreface)
	if n != len(clientPreface) || err != nil {
		return false
	}
	return true
}

func (c *h2c) sendSettingFrame() bool {
	var ss []http2.Setting
	ss = append(ss, http2.Setting{
		ID:  http2.SettingInitialWindowSize,
		Val: 100,
	})

	c.framer.WriteSettings(ss...)
	return true
}

func (c *h2c) sendSettingResponse() bool {
	c.framer.WriteSettings()
	return true
}

func (c *h2c) sendGoAway() bool {
	c.framer.WriteGoAway(c.MaxConcurrentStreams, http2.ErrCodeStreamClosed, nil)
	return true
}

func readResponse(c *h2c) bool {
	for {
		frame, err := c.framer.ReadFrame()
		if err != nil {
			fmt.Println("read err:", err)
			return false
		}
		fmt.Println("recv frame:", frame)
		switch frame.Header().Type {
		case http2.FrameHeaders:
			return true
		case http2.FrameSettings:
		case http2.FrameWindowUpdate:
		case http2.FrameRSTStream:
			fmt.Println("server says rst")
			return true
		case http2.FrameData:
			fmt.Println("recv data frame", frame)
		case http2.FrameGoAway:
			fmt.Println("server says go away")
			return true
		default:
			fmt.Println("unkonwn reponses frame type", frame)
			return true
		}
	}
}

func doTest(c *h2c) {
	streamid := uint32(1)

	c.sendPreface()
	c.sendSettingFrame()

	frame, err := c.framer.ReadFrame()
	if err != nil {
		fmt.Println("readframe", err)
		return
	}
	if frame.Header().Type != http2.FrameSettings {
		fmt.Println("first frame should be SETTINGS")
		return
	}
	settings := frame.(*http2.SettingsFrame)
	settings.ForeachSetting(func(s http2.Setting) error {
		switch s.ID {
		case http2.SettingMaxConcurrentStreams:
			c.MaxConcurrentStreams = s.Val
		}
		return nil
	})

	c.sendSettingResponse()

	fmt.Println("http2 init done, sending requests. 10 RPS.")

	for {
		var headerWriteBuf bytes.Buffer
		enc := hpack.NewEncoder(&headerWriteBuf)
		enc.WriteField(hpack.HeaderField{Name: ":method", Value: "GET"})
		enc.WriteField(hpack.HeaderField{Name: ":scheme", Value: "http"})
		enc.WriteField(hpack.HeaderField{Name: ":path", Value: "/"})
		enc.WriteField(hpack.HeaderField{Name: ":authority", Value: ServerAddr})
		c.framer.WriteHeaders(http2.HeadersFrameParam{
			StreamID:      streamid,
			BlockFragment: headerWriteBuf.Bytes(),
			EndStream:     true,
			EndHeaders:    true,
		})

		if !readResponse(c) {
			return
		}

		streamid += 2
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	fmt.Printf("Connecting to %s\n", ServerAddr)

	conn, _ := net.Dial("tcp", ServerAddr)

	fmt.Println("tcp connected")
	fmt.Println("sleep for 15 seconds, exec \"nginx -s reload\" now")
	time.Sleep(15 * time.Second)
	fmt.Println("sleep done, sending http2 preface")

	/* for tls test case */
	// tlsconn := tls.Client(conn, &tls.Config{
	// 	NextProtos:         []string{"h2"},
	// 	InsecureSkipVerify: true,
	// })
	// tlsconn.Handshake()
	// conn = tlsconn

	c := h2c{
		conn:   conn,
		framer: http2.NewFramer(conn, conn),
	}
	go func() {
		doTest(&c)
		wg.Done()
	}()

	wg.Wait()
}
