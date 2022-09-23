//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package internal

import (
	"fmt"
	"goim-client/api/grpc"
	"net"
	"time"
)

func HandleTCPMsg(c *net.TCPConn, p *grpc.Proto) {
	switch p.Op {
	case 8:
		if !hbOpen {
			go DocSyncTaskCronJobByTCP(c)
			hbOpen = true
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
			fmt.Printf("messageReceivedByWS： ver=%v,body=%s\n", p.Ver, string(body[16:]))
		}
	case 1000:
		fmt.Printf("messageReceivedByWS： ver=%v,body=%s\n", p.Ver, string(p.Body))
	default:
		fmt.Println("未识别的指令", p.Op)
	}
}

func DocSyncTaskCronJobByTCP(c *net.TCPConn) {
	for {
		if heartbeatInterval.After(time.Now()) {
			HeartbeatTCP(c)
			heartbeatInterval.Add(hearBeatSpec)
			time.Sleep(time.Second * 10)
		}
	}
}

func HeartbeatTCP(c *net.TCPConn) {
	msgBuf := Heartbeat()
	c.Write(msgBuf)
}
