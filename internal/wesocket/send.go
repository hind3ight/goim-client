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
	"goim-client/internal"
	"time"
)

// 权限校验

func AuthWS(c *websocket.Conn) {
	msgBuf := internal.Auth()
	_ = c.WriteMessage(websocket.BinaryMessage, msgBuf)
	fmt.Println("send: auth")
}

// 心跳

func HeartbeatWS(c *websocket.Conn) {
	msgBuf := internal.Heartbeat()
	_ = c.WriteMessage(websocket.BinaryMessage, msgBuf)
	fmt.Println("send: heartbeat")
}

// 向服务器端发消息

func SendMsgByWS(ws *websocket.Conn) {
	timer := time.Tick(internal.SendMsgSpec)
	for t := range timer {
		msg := []byte("sendMsg by ws client,data: " + t.Format("2006-01-02 15:04:05"))
		_ = ws.WriteMessage(websocket.BinaryMessage, internal.PackageMsg(msg))
		fmt.Printf("sendMsgByWS :%s\n", string(msg))

	}
}
