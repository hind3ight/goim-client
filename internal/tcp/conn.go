//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package tcp

import (
	"fmt"
	"goim-client/internal"
	"log"
	"net"
	"time"
)

var (
	reconnectSpec = time.Minute * 1
)

const (
	tcpUrl = `192.168.32.97:3101`
	//tcpUrl = `127.0.0.1:3101`
)

// 创建连接

func CreateTCPConn() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpUrl)
	if err != nil {
		log.Printf("Resolve tcp addr failed: %v\n", err)
		return
	}

	// 向服务器拨号
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		time.Sleep(reconnectSpec)
		reconnect <- struct{}{}
		log.Printf("Dial to server failed: %v\n", err)
		return
	}

	// 向服务器发消息
	AuthTCP(conn)
	go OnMessage(conn)
	go SendMsgByTCP(conn)

}

var reconnect = make(chan struct{})

func Reconnect() {
	for {
		select {
		case <-reconnect:
			fmt.Printf("%s开始重连", time.Now().Format("2006-01-02 15:04:05"))
			CreateTCPConn()         // 创建新的连接
			internal.HbOpen = false // 重置心跳服务

		}
	}
}

func DoHeartCronJob(c *net.TCPConn) {
	timer := time.NewTimer(internal.HearBeatSpec)
	HeartbeatTCP(c)
	for {
		select {
		case <-timer.C:
			HeartbeatTCP(c)
		}
	}
}
