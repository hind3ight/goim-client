//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"goim-client/api/grpc"
	"goim-client/pkg/encoding/binary"
	"strings"
	"time"
)

var hbOpen bool

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
		return p, errors.New("长度错误")
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

func ParseMsg2(msg []byte) (err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)

	p := &grpc.Proto{}
	//buf = msg
	if len(buf) < _rawHeaderSize {
		return errors.New("长度错误")
	}
	packLen = binary.BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Int16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Int32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Int32(buf[_seqOffset:])
	if packLen < 0 || packLen > _maxPackSize {
		return errors.New("")
	}
	if headerLen != _rawHeaderSize {
		return errors.New("")
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
		fmt.Printf("messageReceived: ver=%v,body=%s\n", p.Ver, p.Body)
	} else {
		p.Body = nil
	}

	return
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
			fmt.Printf("messageReceivedByHTTP： ver=%v,body=%s\n", p.Ver, string(body[16:]))
		}
	case 5:
		fmt.Printf("messageReceivedByWS: ver=%v,body=%s\n", p.Ver, p.Body)
	default:
		fmt.Println("未识别的指令", p.Op)
	}
}

// 权限校验
func Auth(c *websocket.Conn) {
	token := TokenStruct{ // todo 根据配置文件读取
		Mid:      125,
		RoomId:   "live://1000",
		Platform: "web",
		Accepts:  []int{1000, 1001, 1002},
	}
	headerBuf := make([]byte, 16)
	bodyBuf := handleJson(token)

	binary.BigEndian.PutInt32(headerBuf[_packOffset:], int32(len(bodyBuf)+_rawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[_verOffset:], 1)
	binary.BigEndian.PutInt32(headerBuf[_opOffset:], 7)
	binary.BigEndian.PutInt32(headerBuf[_seqOffset:], 1)

	c.WriteMessage(websocket.BinaryMessage, mergeArrayBuffer(headerBuf, bodyBuf))
}
func mergeArrayBuffer(headerBuf, bodyBuf []byte) (res []byte) {
	res = append(headerBuf, bodyBuf...)
	return
}

func handleJson(token TokenStruct) []byte {
	b, _ := json.Marshal(token)
	tmpSlice := strings.Split(string(b), `,"`)
	b = []byte(strings.Join(tmpSlice, `, "`))
	return b
}

// 心跳
func Heartbeat(ws *websocket.Conn) {
	heartBeatBuf := make([]byte, 16)
	binary.BigEndian.PutInt32(heartBeatBuf[_packOffset:], int32(_rawHeaderSize))
	binary.BigEndian.PutInt16(heartBeatBuf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(heartBeatBuf[_verOffset:], 1)
	binary.BigEndian.PutInt32(heartBeatBuf[_opOffset:], 2)
	binary.BigEndian.PutInt32(heartBeatBuf[_seqOffset:], 1)
	err := ws.WriteMessage(websocket.BinaryMessage, heartBeatBuf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("send: heartbeat")
}

func SendMsgByWS(ws *websocket.Conn, msg []byte) {
	headerBuf := make([]byte, 16)

	binary.BigEndian.PutInt32(headerBuf[_packOffset:], int32(len(msg)+_rawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[_verOffset:], 1)
	binary.BigEndian.PutInt32(headerBuf[_opOffset:], 4)
	binary.BigEndian.PutInt32(headerBuf[_seqOffset:], 1)

	ws.WriteMessage(websocket.BinaryMessage, mergeArrayBuffer(headerBuf, msg))

	return
}

// 定时任务
var heartbeatInterval time.Time

const hearBeatSpec = time.Second * 30

func DocSyncTaskCronJob(c *websocket.Conn) {
	for {
		if heartbeatInterval.After(time.Now()) {
			Heartbeat(c)
			heartbeatInterval.Add(hearBeatSpec)
			time.Sleep(time.Second * 10)
		}
	}
}
