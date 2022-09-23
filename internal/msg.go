//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package internal

import (
	"encoding/json"
	"errors"
	"github.com/Terry-Mao/goim/pkg/encoding/binary"
	"goim-client/api/grpc"
	"io"
	"strings"
)

var token = TokenStruct{ // todo 根据配置文件读取
	Mid:      1242,
	RoomId:   "live://1000",
	Platform: "web",
	Accepts:  []int{1000, 1001, 1002},
}

// 权限
func Auth() (msg []byte) {

	headerBuf := make([]byte, 16)
	bodyBuf := handleJson(token)

	binary.BigEndian.PutInt32(headerBuf[PackOffset:], int32(len(bodyBuf)+RawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[HeaderOffset:], int16(RawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[VerOffset:], 1)
	binary.BigEndian.PutInt32(headerBuf[OpOffset:], 7)
	binary.BigEndian.PutInt32(headerBuf[SeqOffset:], 1)
	return mergeArrayBuffer(headerBuf, bodyBuf)

}

// 心跳
func Heartbeat() []byte {
	heartBeatBuf := make([]byte, 16)
	binary.BigEndian.PutInt32(heartBeatBuf[PackOffset:], int32(RawHeaderSize))
	binary.BigEndian.PutInt16(heartBeatBuf[HeaderOffset:], int16(RawHeaderSize))
	binary.BigEndian.PutInt16(heartBeatBuf[VerOffset:], 1)
	binary.BigEndian.PutInt32(heartBeatBuf[OpOffset:], 2)
	binary.BigEndian.PutInt32(heartBeatBuf[SeqOffset:], 1)
	return heartBeatBuf
}

func handleJson(token TokenStruct) []byte {
	b, _ := json.Marshal(token)
	tmpSlice := strings.Split(string(b), `,"`)
	b = []byte(strings.Join(tmpSlice, `, "`))
	return b
}

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
	if len(buf) < RawHeaderSize {
		return p, io.EOF
	}
	packLen = binary.BigEndian.Int32(buf[PackOffset:HeaderOffset])
	headerLen = binary.BigEndian.Int16(buf[HeaderOffset:VerOffset])
	p.Ver = int32(binary.BigEndian.Int16(buf[VerOffset:OpOffset]))
	p.Op = binary.BigEndian.Int32(buf[OpOffset:SeqOffset])
	p.Seq = binary.BigEndian.Int32(buf[SeqOffset:])
	if packLen < 0 || packLen > _maxPackSize {
		return p, errors.New("")
	}
	if headerLen != RawHeaderSize {
		return p, errors.New("")
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		p.Body = buf[headerLen:packLen]
	} else {
		p.Body = nil
	}

	return p, nil
}

// 合并[]byte
func mergeArrayBuffer(headerBuf, bodyBuf []byte) (res []byte) {
	res = append(headerBuf, bodyBuf...)
	return
}

// 向comet发送的信息
func PackageMsg(msg []byte) []byte {
	headerBuf := make([]byte, 16)

	binary.BigEndian.PutInt32(headerBuf[PackOffset:], int32(len(msg)+RawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[HeaderOffset:], int16(RawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[VerOffset:], 1)
	binary.BigEndian.PutInt32(headerBuf[OpOffset:], 4)
	binary.BigEndian.PutInt32(headerBuf[SeqOffset:], 1)
	return mergeArrayBuffer(headerBuf, msg)
}
