package tplink

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"time"
)

const _timeOut = 5 * time.Second

type Conn interface {
	Close() error
	SetDeadline(time.Time) error
	Write([]byte) (int, error)
	Read([]byte) (int, error)
}

type ConnectionFactory func(string, string, time.Duration) (Conn, error)

var TCPConnFactory = func(proto, addr string, t time.Duration) (Conn, error) {
	return net.DialTimeout(proto, addr, t)
}

func command(cf ConnectionFactory, addr string, cmd interface{}) ([]byte, error) {
	payload, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	conn, err := cf("tcp", addr, _timeOut)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(_timeOut)); err != nil {
		return nil, err
	}
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(payload)))
	bs := append(header, autokeyeEncrypt(payload)...)
	_, err = conn.Write(bs)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Read(header); err != nil {
		return nil, err
	}
	buf := make([]byte, 40*1024)
	l, rErr := conn.Read(buf)
	if rErr != nil && rErr != io.EOF {
		return nil, rErr
	}
	return autokeyeDecrypt(buf[0:l]), nil
}
