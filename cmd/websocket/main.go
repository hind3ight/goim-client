//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	ws "goim-client/internal/websocket"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "192.168.32.124:3102", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/sub"}
	log.Printf("connecting to %s", u.String())
	fmt.Println(u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	authMsg := ws.Auth()
	err = c.WriteMessage(websocket.BinaryMessage, authMsg)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			p, err := ws.ParseMsg(message)
			if err != nil {
				log.Println("read:", err)
				return
			}
			ws.HandleMsg(c, p)
		}
	}()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return

		case t := <-ticker.C:
			type msgStruct struct {
				Time string `json:"time"`
			}
			msg := msgStruct{Time: t.String()}
			msgByte, _ := json.Marshal(msg)
			//ws.SendMsgByWS(c, msgByte)

			ws.SendMsgByHttp(randType(), msgByte)
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

	// todo 重新建立连接
}

// 随机
func randType() (sendType int) {
lib:
	rand.Seed(time.Now().UnixNano())
	sendType = rand.Intn(3)
	if sendType == 0 {
		goto lib
	}
	return
}
