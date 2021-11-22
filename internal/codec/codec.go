package codec

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/gim/internal/codec/pb"
	"github.com/yddeng/utils/buffer"
	"io"
	"reflect"
)

var ErrTooLarge = fmt.Errorf("Message too large")

const (
	seqSize  = 4
	cmdSize  = 2
	bodySize = 2
	HeadSize = seqSize + cmdSize + bodySize
	BuffSize = 65535
)

type Codec struct {
	ns       string
	readBuf  *buffer.Buffer
	readHead bool
	seqNo    uint32
	cmd      uint16
	bodyLen  uint16
}

func NewCodec(ns string) *Codec {
	return &Codec{
		ns:      ns,
		readBuf: buffer.NewBufferWithCap(BuffSize),
	}
}

//解码
func (decoder *Codec) Decode(reader io.Reader) (interface{}, error) {
	for {
		msg, err := decoder.unPack()

		if msg != nil {
			return msg, nil

		} else if err == nil {
			_, err1 := decoder.readBuf.ReadFrom(reader)
			if err1 != nil {
				return nil, err1
			}
		} else {
			return nil, err
		}
	}
}

func (decoder *Codec) unPack() (interface{}, error) {

	if !decoder.readHead {
		if decoder.readBuf.Len() < HeadSize {
			return nil, nil
		}

		decoder.seqNo, _ = decoder.readBuf.ReadUint32BE()
		decoder.cmd, _ = decoder.readBuf.ReadUint16BE()
		decoder.bodyLen, _ = decoder.readBuf.ReadUint16BE()
		decoder.readHead = true
	}

	if decoder.bodyLen > BuffSize-HeadSize {
		return nil, ErrTooLarge
	}
	if decoder.readBuf.Len() < int(decoder.bodyLen) {
		return nil, nil
	}

	data, _ := decoder.readBuf.ReadBytes(int(decoder.bodyLen))
	m, err := pb.Unmarshal(decoder.ns, decoder.cmd, data)
	if err != nil {
		return nil, err
	}
	msg := &Message{
		data: m.(proto.Message),
		cmd:  decoder.cmd,
		seq:  decoder.seqNo,
	}

	decoder.readHead = false
	return msg, nil
}

//编码
func (encoder *Codec) Encode(o interface{}) ([]byte, error) {
	var seqNo uint32
	var cmd uint16
	var data []byte
	var bodyLen int
	var err error

	switch o.(type) {
	case *Message:
		msg := o.(*Message)
		cmd, data, err = pb.Marshal(encoder.ns, msg.GetData())
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid type:%s", reflect.TypeOf(o).String())
	}

	bodyLen = len(data)
	if bodyLen > BuffSize-HeadSize {
		return nil, ErrTooLarge
	}

	totalLen := HeadSize + bodyLen
	buff := buffer.NewBufferWithCap(totalLen)
	//seq
	buff.WriteUint32BE(seqNo)
	//cmd
	buff.WriteUint16BE(cmd)
	//bodylen
	buff.WriteUint16BE(uint16(bodyLen))
	//body
	buff.WriteBytes(data)

	return buff.Bytes(), nil
}
