//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package internal

import "time"

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)
const (
	// size
	PackSize      = 4
	HeaderSize    = 2
	VerSize       = 2
	OpSize        = 4
	SeqSize       = 4
	HeartSize     = 4
	RawHeaderSize = PackSize + HeaderSize + VerSize + OpSize + SeqSize
	_maxPackSize  = MaxBodySize + int32(RawHeaderSize)
	// offset
	PackOffset   = 0
	HeaderOffset = PackOffset + PackSize
	VerOffset    = HeaderOffset + HeaderSize
	OpOffset     = VerOffset + VerSize
	SeqOffset    = OpOffset + OpSize
	HeartOffset  = SeqOffset + SeqSize
)

var HbOpen bool

// 定时任务
var HeartbeatInterval time.Time

const HearBeatSpec = time.Second * 30 // 心跳间隔时间
const SendMsgSpec = time.Second * 10  // 消息发送间隔时间
