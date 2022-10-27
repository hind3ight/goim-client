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
	"goim-client/server"
	"io"
	"log"
	"net"
)

func OnMessage(c *net.TCPConn) {
	// 接收来自服务器端的广播消息
	buf := make([]byte, 1024)
	for {
		length, err := c.Read(buf)
		realBuf := buf[:length]
		p, err := internal.ParseMsg(realBuf)
		if err != nil {
			if err == io.EOF {
				reconnect <- struct{}{}
				c.Close()
				return
			}
			log.Println("读取错误:", err)
			return
		}
		handleTCPMsg(c, p)
	}
}

func handleTCPMsg(c *net.TCPConn, p *internal.MsgProto) {
	switch p.Op {
	case 8:
		if !internal.HbOpen {
			go DoHeartCronJob(c)
			internal.HbOpen = true
		}
	case 3:
		fmt.Println("receive: heartbeat")
	case 9:

		if len(p.Body) > 16 {
			fmt.Printf("messageReceived： ver=%v,body=%s\n", p.Ver, string(p.Body[16:]))
		}
	case 5:

		if len(p.Body) > 16 {
			fmt.Printf("messageReceived： ver=%v,body=%s\n", p.Ver, string(p.Body[16:]))
		}
	case 1000:
		fmt.Printf("messageReceived： ver=%v,body=%s\n", p.Ver, string(p.Body))
	default:
		fmt.Println("未识别的指令", p.Op)
	}
}
