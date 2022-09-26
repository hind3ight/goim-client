//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package tcp

import (
	"fmt"
	"goim-client/api/grpc"
	"goim-client/internal"
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
				return
			}
			log.Println("读取错误:", err)
			return
		}
		handleTCPMsg(c, p)
	}
}

func handleTCPMsg(c *net.TCPConn, p *grpc.Proto) {
	switch p.Op {
	case 8:
		if !internal.HbOpen {
			go DoHeartCronJob(c)
			internal.HbOpen = true
		}
	case 3:
		fmt.Println("receive: heartbeat")
	case 9:
		body := p.GetBody()

		if len(body) > 16 {
			fmt.Printf("messageReceived： ver=%v,body=%s\n", p.Ver, string(body[16:]))
		}
	case 5:
		body := p.GetBody()

		if len(body) > 16 {
			fmt.Printf("messageReceived： ver=%v,body=%s\n", p.Ver, string(body[16:]))
		}
	case 1000:
		fmt.Printf("messageReceived： ver=%v,body=%s\n", p.Ver, string(p.Body))
	default:
		fmt.Println("未识别的指令", p.Op)
	}
}
