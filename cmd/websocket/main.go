//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	ws "goim-client/internal/websocket"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var addr = flag.String("addr", "192.168.32.124:3102", "http service address")

func main() {

	// signal
	c := make(chan os.Signal, 1)
	createWSConn(c)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		fmt.Printf("get a signal %s", s.String())
		//log.Info()
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			//svc.Close()
			fmt.Println(`service exit`)
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

}

var (
	MaxConnectTimes = 10
	Delay           = 15000
)

func createWSConn(signal chan os.Signal) *websocket.Conn {

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/sub"}
	log.Printf("connecting to %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	ws.Auth(conn)
	ws.SendMsgByWS(conn, []byte("hello"))
	go onMessage(conn)
	tick := time.Tick(time.Second * 10)
	go func() {
		for {
			select {
			case _ = <-tick: // 每隔10秒发送信息
				//fmt.Println(t)
				//ws.SendMsgByWS(conn, []byte(fmt.Sprintf("hello,this is %s", t.Format("2006-01-02 15:04:05"))))
			case <-signal:
				log.Println("interrupt")

				err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}

				return

			}
		}
	}()
	return conn
}

func connect(u url.URL) {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	go onMessage(c)
}

func onMessage(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()

		p, err := ws.ParseMsg(message)
		if err != nil {
			log.Println("read:", err)
			return
		}
		ws.HandleMsg(c, p)
	}
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
