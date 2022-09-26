//
//	@Description
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package internal

type TokenStruct struct {
	Mid      int    `json:"mid"`
	RoomId   string `json:"room_id"`
	Platform string `json:"platform"`
	Accepts  []int  `json:"accepts"`
}

type MsgProto struct {
	Ver  int32  `protobuf:"varint,1,opt,name=ver,proto3" json:"ver"`
	Op   int32  `protobuf:"varint,2,opt,name=op,proto3" json:"op"`
	Seq  int32  `protobuf:"varint,3,opt,name=seq,proto3" json:"seq"`
	Body []byte `protobuf:"bytes,4,opt,name=body,proto3" json:"body"`
}
