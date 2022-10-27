//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package wesocket

import (
	"fmt"
	"goim-client/internal"
	"io"
	"log"
	"time"
)

// 信息处理
func (s *WSConn) OnMessage() {
	for {
		_, message, err := s.conn.ReadMessage()

		p, err := internal.ParseMsg(message)
		if err != nil {
			if err == io.EOF {
				reconnect <- struct{}{}
				connWs <- s.conn
				s.conn.Close()
				return
			}
			log.Println("读取错误:", err)
			return
		}
		s.handleWSMsg(p)
	}
}

// 根据msg处理

func (s *WSConn) handleWSMsg(p *internal.MsgProto) {
	switch p.Op {
	case 8:
		if !internal.HbOpen {
			timer := time.NewTimer(internal.HearBeatSpec)
			go s.DoHeartCronJob(timer)
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
	default:
		fmt.Println("未识别的指令", p.Op)
	}
}
