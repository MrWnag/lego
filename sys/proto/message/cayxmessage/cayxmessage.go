package cayxmessage

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	legoproto "lego/sys/proto"

	"github.com/golang/protobuf/proto"
)

const tcpinfosize = 4
const msgheadsize = 8
const DK_MAPPED_NEW = 0x08 //校验类型

type CayxMessage struct {
	cbDataKind  byte   //数据类型
	cbCheckCode byte   //效验字段
	wPacketSize uint16 //数据大小
	main_id     uint16 //主Id
	sub_id      uint16 //子id
	buffer      []byte //缓存
}

func (this *CayxMessage) GetComId() uint16 {
	return this.main_id
}

func (this *CayxMessage) GetMsgId() uint16 {
	return this.sub_id
}

func (this *CayxMessage) GetMsg() []byte {
	return this.buffer
}

func (this *CayxMessage) Serializable() (bytes []byte, err error) {
	buffer := make([]byte, 0)
	buffer = append(buffer, this.cbDataKind)
	buffer = append(buffer, this.cbCheckCode)
	if legoproto.IsUseBigEndian {
		tmpbbuffer := make([]byte, 2)
		binary.BigEndian.PutUint16(tmpbbuffer, this.wPacketSize)
		buffer = append(buffer, tmpbbuffer...)
		tmpbbuffer = make([]byte, 2)
		binary.BigEndian.PutUint16(tmpbbuffer, this.main_id)
		buffer = append(buffer, tmpbbuffer...)
		tmpbbuffer = make([]byte, 2)
		binary.BigEndian.PutUint16(tmpbbuffer, this.sub_id)
		buffer = append(buffer, tmpbbuffer...)
		buffer = append(buffer, this.buffer...)
	} else {
		tmpbbuffer := make([]byte, 2)
		binary.LittleEndian.PutUint16(tmpbbuffer, this.wPacketSize)
		buffer = append(buffer, tmpbbuffer...)
		tmpbbuffer = make([]byte, 2)
		binary.LittleEndian.PutUint16(tmpbbuffer, this.main_id)
		buffer = append(buffer, tmpbbuffer...)
		tmpbbuffer = make([]byte, 2)
		binary.LittleEndian.PutUint16(tmpbbuffer, this.sub_id)
		buffer = append(buffer, tmpbbuffer...)
		buffer = append(buffer, this.buffer...)
	}
	mappedBuffer(buffer, len(buffer))
	return buffer, nil
}

func (this *CayxMessage) ToString() string {
	if legoproto.MsgProtoType == legoproto.Proto_Json {
		return fmt.Sprintf("ComId:%d MsgId:%d Msg:%s", this.main_id, this.sub_id, string(this.buffer))
	} else {
		return fmt.Sprintf("ComId:%d MsgId:%d Msg:%v", this.main_id, this.sub_id, this.buffer)
	}
}

func MessageDecodeBybufio(r *bufio.Reader) (msg legoproto.IMessage, err error) {
	cbDataKind, err := legoproto.ReadByte(r)
	if err != nil {
		return nil, err
	}
	if (cbDataKind & DK_MAPPED_NEW) != cbDataKind {
		return nil, fmt.Errorf("消息校验失败 CheckId:0x%x", cbDataKind)
	}
	cbCheckCode, err := legoproto.ReadByte(r)
	if err != nil {
		return nil, err
	}
	wPacketSize, err := legoproto.ReadUInt16(r)
	if err != nil {
		return nil, err
	}
	_tmpbuffer := make([]byte, wPacketSize-tcpinfosize)
	_, err = io.ReadFull(r, _tmpbuffer)

	buffer := make([]byte, 0)
	buffer = append(buffer, cbDataKind)
	buffer = append(buffer, cbCheckCode)
	_msg := make([]byte, 2)
	if legoproto.IsUseBigEndian {
		binary.BigEndian.PutUint16(_msg, wPacketSize)
	} else {
		binary.LittleEndian.PutUint16(_msg, wPacketSize)
	}
	buffer = append(buffer, _msg...)
	buffer = append(buffer, _tmpbuffer...)

	if !unMappedBuffer(buffer, len(buffer)) {
		return nil, fmt.Errorf("解码包体错误")
	}
	message := &CayxMessage{}
	message.cbDataKind = buffer[0]
	message.cbCheckCode = buffer[1]
	if legoproto.IsUseBigEndian {
		message.wPacketSize = binary.BigEndian.Uint16(buffer[2:4])
		message.main_id = binary.BigEndian.Uint16(buffer[4:6])
		message.sub_id = binary.BigEndian.Uint16(buffer[6:8])
	} else {
		message.wPacketSize = binary.LittleEndian.Uint16(buffer[2:4])
		message.main_id = binary.LittleEndian.Uint16(buffer[4:6])
		message.sub_id = binary.LittleEndian.Uint16(buffer[6:8])
	}
	message.buffer = buffer[msgheadsize:]
	return message, err
}

func MessageDecodeBybytes(buffer []byte) (msg legoproto.IMessage, err error) {
	if len(buffer) < msgheadsize {
		return nil, fmt.Errorf("解析数据失败 buffer 长度:%d", len(buffer))
	}
	if !unMappedBuffer(buffer, len(buffer)) {
		return nil, fmt.Errorf("解码包体错误")
	}
	message := &CayxMessage{}
	message.cbDataKind = buffer[0]
	message.cbCheckCode = buffer[1]
	if legoproto.IsUseBigEndian {
		message.wPacketSize = binary.BigEndian.Uint16(buffer[2:4])
		message.main_id = binary.BigEndian.Uint16(buffer[4:6])
		message.sub_id = binary.BigEndian.Uint16(buffer[6:8])
	} else {
		message.wPacketSize = binary.LittleEndian.Uint16(buffer[2:4])
		message.main_id = binary.LittleEndian.Uint16(buffer[4:6])
		message.sub_id = binary.LittleEndian.Uint16(buffer[6:8])
	}
	message.buffer = buffer[msgheadsize:]
	return message, err
}

func MessageMarshal(comId uint16, msgId uint16, msg interface{}) legoproto.IMessage {
	if legoproto.MsgProtoType == legoproto.Proto_Json {
		return jsonDefMessageMarshal(comId, msgId, msg)
	} else {
		return protoDefMessageMarshal(comId, msgId, msg)
	}
}

// Proto 消息 编码  方法
func protoDefMessageMarshal(comId uint16, msgId uint16, msg interface{}) legoproto.IMessage {
	data, _ := proto.Marshal(msg.(proto.Message))
	message := &CayxMessage{
		main_id:     uint16(comId),
		sub_id:      uint16(msgId),
		cbCheckCode: DK_MAPPED_NEW,
		wPacketSize: uint16(len(data) + msgheadsize),
		buffer:      data,
	}
	return message
}

// Json 消息 编码  方法 这个有点坑爹啊 _Msg里面的字段必须是大写开头的哈 不然无法序列化的
func jsonDefMessageMarshal(comId uint16, msgId uint16, msg interface{}) legoproto.IMessage {
	data, _ := json.Marshal(msg)
	message := &CayxMessage{
		main_id:     uint16(comId),
		sub_id:      uint16(msgId),
		cbCheckCode: DK_MAPPED_NEW,
		wPacketSize: uint16(len(data) + msgheadsize),
		buffer:      data,
	}
	return message
}

var ENCODE_MAP_NEW = []byte{
	0x9b, 0x38, 0x32, 0x00, 0x45, 0x83, 0x48, 0xb8, 0xc9, 0xd8, 0xb5, 0x73, 0x13, 0xbe, 0x4c, 0xf4,
	0xfe, 0xe1, 0x42, 0xec, 0x39, 0x65, 0x0e, 0x3f, 0x6e, 0x99, 0xaf, 0xd2, 0x51, 0xf2, 0x3b, 0xa9,
	0xd6, 0x24, 0xe6, 0x27, 0x5c, 0xc2, 0x72, 0xd4, 0x6d, 0xaa, 0xbc, 0xda, 0x03, 0x88, 0x74, 0x0c,
	0xbb, 0x15, 0x91, 0xb6, 0xac, 0xe8, 0xef, 0xeb, 0x0d, 0x22, 0x9c, 0xe7, 0x75, 0x1f, 0x4e, 0xde,
	0x58, 0xdc, 0x67, 0x90, 0xf1, 0x3c, 0x7c, 0x35, 0x59, 0xb0, 0xb2, 0xa0, 0x09, 0xdb, 0xd1, 0x82,
	0xd0, 0xba, 0x1b, 0xdf, 0x53, 0x3a, 0x08, 0x2b, 0xb4, 0x9e, 0x02, 0xd3, 0x68, 0x6a, 0xe9, 0x57,
	0x43, 0x81, 0x46, 0x04, 0xc1, 0xa7, 0x31, 0x0a, 0x98, 0x16, 0x86, 0x3d, 0xff, 0xab, 0x4b, 0xa4,
	0x49, 0x61, 0x7e, 0xf9, 0x33, 0x29, 0xb3, 0x60, 0x1d, 0x9f, 0x8a, 0xc7, 0x6f, 0xf6, 0x8d, 0x26,
	0xd7, 0xee, 0x34, 0xae, 0xa2, 0x63, 0xe2, 0xa8, 0x87, 0xcd, 0x6b, 0xb7, 0x07, 0xf8, 0xbf, 0xc6,
	0x37, 0x06, 0x2e, 0x14, 0xcc, 0x36, 0xa3, 0x7b, 0xad, 0x7d, 0x85, 0x78, 0x79, 0xc4, 0x3e, 0x5b,
	0xca, 0x4d, 0x12, 0xea, 0x6c, 0x41, 0xc3, 0x76, 0x05, 0x77, 0x18, 0x93, 0x44, 0x55, 0x8e, 0xe5,
	0x5e, 0x69, 0x20, 0x89, 0x97, 0x5a, 0x92, 0x56, 0x52, 0x19, 0x70, 0xf5, 0x2c, 0x01, 0xed, 0x2a,
	0x2d, 0x17, 0x71, 0x50, 0x40, 0x0b, 0x7f, 0x1c, 0x84, 0x64, 0x95, 0x21, 0xb9, 0x8f, 0xf0, 0xd5,
	0xfd, 0x0f, 0x80, 0x11, 0x25, 0x1e, 0xe0, 0xfa, 0x94, 0x96, 0xdd, 0x2f, 0x1a, 0xf7, 0xe3, 0xe4,
	0xa6, 0xb1, 0xa5, 0xcf, 0x5f, 0xa1, 0xcb, 0x10, 0x62, 0x9a, 0x54, 0x9d, 0x5d, 0x28, 0xfb, 0x7a,
	0xce, 0x30, 0x47, 0x66, 0xf3, 0xc8, 0x4a, 0xc0, 0x4f, 0xd9, 0x23, 0xbd, 0x8c, 0xfc, 0x8b, 0xc5,
}

var DECODE_MAP_NEW = []byte{
	0x03, 0xbd, 0x5a, 0x2c, 0x63, 0xa8, 0x91, 0x8c, 0x56, 0x4c, 0x67, 0xc5, 0x2f, 0x38, 0x16, 0xd1,
	0xe7, 0xd3, 0xa2, 0x0c, 0x93, 0x31, 0x69, 0xc1, 0xaa, 0xb9, 0xdc, 0x52, 0xc7, 0x78, 0xd5, 0x3d,
	0xb2, 0xcb, 0x39, 0xfa, 0x21, 0xd4, 0x7f, 0x23, 0xed, 0x75, 0xbf, 0x57, 0xbc, 0xc0, 0x92, 0xdb,
	0xf1, 0x66, 0x02, 0x74, 0x82, 0x47, 0x95, 0x90, 0x01, 0x14, 0x55, 0x1e, 0x45, 0x6b, 0x9e, 0x17,
	0xc4, 0xa5, 0x12, 0x60, 0xac, 0x04, 0x62, 0xf2, 0x06, 0x70, 0xf6, 0x6e, 0x0e, 0xa1, 0x3e, 0xf8,
	0xc3, 0x1c, 0xb8, 0x54, 0xea, 0xad, 0xb7, 0x5f, 0x40, 0x48, 0xb5, 0x9f, 0x24, 0xec, 0xb0, 0xe4,
	0x77, 0x71, 0xe8, 0x85, 0xc9, 0x15, 0xf3, 0x42, 0x5c, 0xb1, 0x5d, 0x8a, 0xa4, 0x28, 0x18, 0x7c,
	0xba, 0xc2, 0x26, 0x0b, 0x2e, 0x3c, 0xa7, 0xa9, 0x9b, 0x9c, 0xef, 0x97, 0x46, 0x99, 0x72, 0xc6,
	0xd2, 0x61, 0x4f, 0x05, 0xc8, 0x9a, 0x6a, 0x88, 0x2d, 0xb3, 0x7a, 0xfe, 0xfc, 0x7e, 0xae, 0xcd,
	0x43, 0x32, 0xb6, 0xab, 0xd8, 0xca, 0xd9, 0xb4, 0x68, 0x19, 0xe9, 0x00, 0x3a, 0xeb, 0x59, 0x79,
	0x4b, 0xe5, 0x84, 0x96, 0x6f, 0xe2, 0xe0, 0x65, 0x87, 0x1f, 0x29, 0x6d, 0x34, 0x98, 0x83, 0x1a,
	0x49, 0xe1, 0x4a, 0x76, 0x58, 0x0a, 0x33, 0x8b, 0x07, 0xcc, 0x51, 0x30, 0x2a, 0xfb, 0x0d, 0x8e,
	0xf7, 0x64, 0x25, 0xa6, 0x9d, 0xff, 0x8f, 0x7b, 0xf5, 0x08, 0xa0, 0xe6, 0x94, 0x89, 0xf0, 0xe3,
	0x50, 0x4e, 0x1b, 0x5b, 0x27, 0xcf, 0x20, 0x80, 0x09, 0xf9, 0x2b, 0x4d, 0x41, 0xda, 0x3f, 0x53,
	0xd6, 0x11, 0x86, 0xde, 0xdf, 0xaf, 0x22, 0x3b, 0x35, 0x5e, 0xa3, 0x37, 0x13, 0xbe, 0x81, 0x36,
	0xce, 0x44, 0x1d, 0xf4, 0x0f, 0xbb, 0x7d, 0xdd, 0x8d, 0x73, 0xd7, 0xee, 0xfd, 0xd0, 0x10, 0x6c,
}

func mappedBuffer(buffer []byte, wDataSize int) bool {
	cbCheckCode := 0
	for i := tcpinfosize; i < wDataSize; i++ {
		cbCheckCode = cbCheckCode + int(buffer[i])
		buffer[i] = ENCODE_MAP_NEW[buffer[i]] //update
	}
	cbCheckCode = (^cbCheckCode + 1)
	buffer[0] = DK_MAPPED_NEW
	buffer[1] = byte(cbCheckCode)
	return true
}

func unMappedBuffer(buffer []byte, size int) bool {
	if (buffer[0] & DK_MAPPED_NEW) != 0 {
		cbCheckCode := buffer[1]
		for i := tcpinfosize; i < size; i++ {
			cbCheckCode += DECODE_MAP_NEW[buffer[i]]
			buffer[i] = DECODE_MAP_NEW[buffer[i]]
		}
		//效验
		if cbCheckCode != 0 {
			return false
		}
	}
	return true
}
