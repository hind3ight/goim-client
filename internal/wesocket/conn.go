//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package wesocket

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"goim-client/internal"
	"log"
	"net/url"
	"time"
)

var addr = flag.String("addr", "192.168.32.124:3102", "http service address")
var (
	reconnectSpec = time.Minute * 1
)

// 创建连接

func CreateWSConn() {

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/sub"}
	log.Printf("connecting to %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		time.Sleep(reconnectSpec)
		reconnect <- struct{}{}
		log.Fatal("dial:", err)
		return
	}
	AuthWS(conn)

	go SendMsgByWS(conn)
	go OnMessage(conn)

	return
}

var reconnect = make(chan struct{})

func Reconnect() {
	for {
		select {
		case <-reconnect:
			fmt.Printf("%s开始重连", time.Now().Format("2006-01-02 15:04:05"))
			CreateWSConn()          // 创建新的连接
			internal.HbOpen = false // 重置心跳服务

		}
	}
}

// 心跳定时任务
func DoHeartCronJob(c *websocket.Conn, timer *time.Timer) {
	HeartbeatWS(c) // 第一次直接执行
	for {
		select {
		case <-timer.C:
			HeartbeatWS(c)
		}
	}
}