//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package wesocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"goim-client/api/grpc"
	"goim-client/internal"
	"io"
	"log"
	"time"
)

// 信息处理
func OnMessage(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()

		p, err := internal.ParseMsg(message)
		if err != nil {
			if err == io.EOF {
				reconnect <- struct{}{}
				c.Close()
				return
			}
			log.Println("读取错误:", err)
			return
		}
		handleWSMsg(c, p)
	}
}

// 根据msg处理

func handleWSMsg(c *websocket.Conn, p *grpc.Proto) {
	switch p.Op {
	case 8:
		if !internal.HbOpen {
			timer := time.NewTimer(internal.HearBeatSpec)
			go DoHeartCronJob(c, timer)
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
	default:
		fmt.Println("未识别的指令", p.Op)
	}
}
