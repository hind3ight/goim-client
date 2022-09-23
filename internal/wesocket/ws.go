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
	"goim-client/pkg/encoding/binary"
	"time"
)

// 根据msg处理
func HandleMsg(c *websocket.Conn, p *grpc.Proto) {
	switch p.Op {
	case 8:
		if !internal.HbOpen {
			go DocSyncTaskCronJob(c)
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
			fmt.Printf("messageReceivedByWS： ver=%v,body=%s\n", p.Ver, string(body[16:]))
		}
	default:
		fmt.Println("未识别的指令", p.Op)
	}
}

// 权限校验

func AuthWS(c *websocket.Conn) {
	msgBuf := internal.Auth()
	c.WriteMessage(websocket.BinaryMessage, msgBuf)
}

// 心跳

func HeartbeatWS(c *websocket.Conn) {
	msgBuf := internal.Heartbeat()
	c.WriteMessage(websocket.BinaryMessage, msgBuf)
}

func SendMsgByWS(ws *websocket.Conn, msg []byte) {
	headerBuf := make([]byte, 16)

	binary.BigEndian.PutInt32(headerBuf[internal.PackOffset:], int32(len(msg)+internal.RawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[internal.HeaderOffset:], int16(internal.RawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[internal.VerOffset:], 1)
	binary.BigEndian.PutInt32(headerBuf[internal.OpOffset:], 4)
	binary.BigEndian.PutInt32(headerBuf[internal.SeqOffset:], 1)

	ws.WriteMessage(websocket.BinaryMessage, internal.PackageMsg(msg))

	return
}

func DocSyncTaskCronJob(c *websocket.Conn) {
	for {
		if internal.HeartbeatInterval.After(time.Now()) {
			HeartbeatWS(c)
			internal.HeartbeatInterval.Add(internal.HearBeatSpec)
			time.Sleep(time.Second * 10)
		}
	}
}
