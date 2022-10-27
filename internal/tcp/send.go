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
	"net"
	"time"
)

// 权限校验

func AuthTCP(c *net.TCPConn) {
	_, _ = c.Write(internal.Auth())
	fmt.Println("send: auth")
}

// 心跳

func HeartbeatTCP(c *net.TCPConn) {
	msgBuf := internal.Heartbeat()
	_, _ = c.Write(msgBuf)
}

// 向服务器端发消息

func SendMsgByTCP(conn net.Conn) {
	timer := time.Tick(internal.SendMsgSpec)
	for t := range timer {
		msg := []byte("sendMsg by tcp client,data: " + t.Format("2006-01-02 15:04:05"))
		if _, err := conn.Write(internal.PackageMsg(msg)); err != nil {
			fmt.Printf("err :%v", err)
			reconnect <- struct{}{}
		} else {
			fmt.Printf("sendMsgByTCP :%s\n", string(msg))
		}
	}
}
