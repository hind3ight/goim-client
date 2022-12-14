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

func (s *WSConn) AuthWS() {
	msgBuf := internal.Auth()
	s.mutex.Lock()
	_ = s.conn.WriteMessage(websocket.BinaryMessage, msgBuf)
	s.mutex.Unlock()

	fmt.Println("send: auth")
}

// 心跳

func (s *WSConn) HeartbeatWS() {
	msgBuf := internal.Heartbeat()
	s.mutex.Lock()
	_ = s.conn.WriteMessage(websocket.BinaryMessage, msgBuf)
	s.mutex.Unlock()

	fmt.Println("send: heartbeat")
}

// 向服务器端发消息

func (s *WSConn) SendMsgByWS() {
	timer := time.Tick(internal.SendMsgSpec)
	for t := range timer {
		select {
		case a := <-connWs:
			if a == s.conn {
				a.Close()
				return
			}
		default:
			msg := []byte("sendMsg by ws client,data: " + t.Format("2006-01-02 15:04:05"))
			s.mutex.Lock()
			err := s.conn.WriteMessage(websocket.BinaryMessage, internal.PackageMsg(msg))
			s.mutex.Unlock()
			if err != nil {
				fmt.Println(err)
				s.conn.Close()
				return
			} else {
				fmt.Printf("sendMsgByWS :%s\n", string(msg))
			}
		}
	}
}

func (s *WSConn) SendHearBeatWS() {
	timer := time.Tick(internal.HearBeatSpec)
	for range timer {
		select {
		case a := <-connWs:
			if a == s.conn {
				a.Close()
				return
			}
		default:
			s.HeartbeatWS()
		}
	}
}
