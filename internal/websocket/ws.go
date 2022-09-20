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

// 根据msg处理
func HandleMsg(c *websocket.Conn, p *grpc.Proto) {
	switch p.Op {
	case 8:
		time.Sleep(time.Second * 30)
		Heartbeat(c)
	case 3:
		fmt.Println("heartbeat reply")
	case 9:
		body := p.GetBody()

		if len(body) > 16 {
			// todo 处理读取到的信息
			fmt.Printf("接受到信息：%s\n", string(body[16:]))
		}
	default:
		fmt.Println(p.Op)
	}
}

// 权限校验
func Auth() []byte {
	token := TokenStruct{ // todo 根据配置文件读取
		Mid:      123,
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
	return mergeArrayBuffer(headerBuf, bodyBuf)
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
	tmpB := make([]byte, 16)
	binary.BigEndian.PutInt32(tmpB[_packOffset:], int32(_rawHeaderSize))
	binary.BigEndian.PutInt16(tmpB[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(tmpB[_verOffset:], 1)
	binary.BigEndian.PutInt32(tmpB[_opOffset:], 2)
	binary.BigEndian.PutInt32(tmpB[_seqOffset:], 1)
	ws.WriteMessage(websocket.BinaryMessage, tmpB)

	fmt.Println("send: heartbeat")
}

func SendMsgByWS(ws *websocket.Conn, msg []byte) {

	return
}
