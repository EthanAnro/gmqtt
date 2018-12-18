package packets

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Suback struct {
	FixHeader *FixHeader
	PacketID  PacketID
	Payload   []byte
}

func (c *Suback) String() string {
	return fmt.Sprintf("Suback, Pid: %v, Payload: %v", c.PacketID, c.Payload)
}

//new suback
func NewSubackPacket(fh *FixHeader, r io.Reader) (*Suback, error) {
	p := &Suback{FixHeader: fh}
	//判断 标志位 flags 是否合法[MQTT-3.8.1-1]
	if fh.Flags != FLAG_RESERVED {
		return nil, ErrInvalFlags
	}
	err := p.Unpack(r)
	return p, err
}

func (p *Suback) Pack(w io.Writer) error {
	p.FixHeader.Pack(w)
	pid := make([]byte, 2)
	binary.BigEndian.PutUint16(pid, p.PacketID)
	w.Write(pid)
	_, err := w.Write(p.Payload)
	return err
}

func (p *Suback) Unpack(r io.Reader) error {
	restBuffer := make([]byte, p.FixHeader.RemainLength)
	_, err := io.ReadFull(r, restBuffer)
	if err != nil {
		return err
	}
	p.PacketID = binary.BigEndian.Uint16(restBuffer[0:2])
	p.Payload = restBuffer[2:]
	return nil
}
