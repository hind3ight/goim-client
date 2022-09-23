//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package internal

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"goim-client/api/grpc"
	"goim-client/pkg/encoding/binary"
	"io"
	"time"
)

var hbOpen bool

// 定时任务
var heartbeatInterval time.Time

const hearBeatSpec = time.Second * 30 // 心跳间隔时间

// 解析msg
func ParseMsg(msg []byte) (pc *grpc.Proto, err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)

	p := &grpc.Proto{}
	buf = msg
	if len(buf) < _rawHeaderSize {
		return p, io.EOF
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen < 0 || packLen > _maxPackSize {
		return p, errors.New("")
	}
	if headerLen != _rawHeaderSize {
		return p, errors.New("")
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}

	return p, nil
}

// 根据msg处理
func HandleMsg(c *websocket.Conn, p *grpc.Proto) {
	switch p.Op {
	case 8:
		if !hbOpen {
			go DocSyncTaskCronJob(c)
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
	default:
		fmt.Println("未识别的指令", p.Op)
	}
}

// 权限校验

func AuthWS(c *websocket.Conn) {
	msgBuf := Auth()
	c.WriteMessage(websocket.BinaryMessage, msgBuf)
}

func mergeArrayBuffer(headerBuf, bodyBuf []byte) (res []byte) {
	res = append(headerBuf, bodyBuf...)
	return
}

// 心跳

func HeartbeatWS(c *websocket.Conn) {
	msgBuf := Heartbeat()
	c.WriteMessage(websocket.BinaryMessage, msgBuf)
}

func SendMsgByWS(ws *websocket.Conn, msg []byte) {
	//headerBuf := make([]byte, 16)
	//
	//binary.BigEndian.PutInt32(headerBuf[_packOffset:], int32(len(msg)+_rawHeaderSize))
	//binary.BigEndian.PutInt16(headerBuf[_headerOffset:], int16(_rawHeaderSize))
	//binary.BigEndian.PutInt16(headerBuf[_verOffset:], 1)
	//binary.BigEndian.PutInt32(headerBuf[_opOffset:], 4)
	//binary.BigEndian.PutInt32(headerBuf[_seqOffset:], 1)

	ws.WriteMessage(websocket.BinaryMessage, PackageMsg(msg))

	return
}

func PackageMsg(msg []byte) []byte {
	headerBuf := make([]byte, 16)

	binary.BigEndian.PutInt32(headerBuf[_packOffset:], int32(len(msg)+_rawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[_verOffset:], 1)
	binary.BigEndian.PutInt32(headerBuf[_opOffset:], 4)
	binary.BigEndian.PutInt32(headerBuf[_seqOffset:], 1)
	return mergeArrayBuffer(headerBuf, msg)
}

func DocSyncTaskCronJob(c *websocket.Conn) {
	for {
		if heartbeatInterval.After(time.Now()) {
			HeartbeatWS(c)
			heartbeatInterval.Add(hearBeatSpec)
			time.Sleep(time.Second * 10)
		}
	}
}
