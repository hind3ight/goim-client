//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package internal

import (
	"encoding/json"
	"github.com/Terry-Mao/goim/pkg/encoding/binary"
	"strings"
)

var token = TokenStruct{ // todo 根据配置文件读取
	Mid:      1242,
	RoomId:   "live://1000",
	Platform: "web",
	Accepts:  []int{1000, 1001, 1002},
}

func Auth() (msg []byte) {

	headerBuf := make([]byte, 16)
	bodyBuf := handleJson(token)

	binary.BigEndian.PutInt32(headerBuf[_packOffset:], int32(len(bodyBuf)+_rawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(headerBuf[_verOffset:], 1)
	binary.BigEndian.PutInt32(headerBuf[_opOffset:], 7)
	binary.BigEndian.PutInt32(headerBuf[_seqOffset:], 1)
	return mergeArrayBuffer(headerBuf, bodyBuf)

}

func Heartbeat() []byte {
	heartBeatBuf := make([]byte, 16)
	binary.BigEndian.PutInt32(heartBeatBuf[_packOffset:], int32(_rawHeaderSize))
	binary.BigEndian.PutInt16(heartBeatBuf[_headerOffset:], int16(_rawHeaderSize))
	binary.BigEndian.PutInt16(heartBeatBuf[_verOffset:], 1)
	binary.BigEndian.PutInt32(heartBeatBuf[_opOffset:], 2)
	binary.BigEndian.PutInt32(heartBeatBuf[_seqOffset:], 1)
	return heartBeatBuf
}

func handleJson(token TokenStruct) []byte {
	b, _ := json.Marshal(token)
	tmpSlice := strings.Split(string(b), `,"`)
	b = []byte(strings.Join(tmpSlice, `, "`))
	return b
}
