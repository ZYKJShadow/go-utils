package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"net"
)

type TCPClient struct {
	conn net.Conn
	r    *bufio.Reader
	w    *bufio.Writer
	head int32
	tag  int32
}

func NewTCPClient(conn net.Conn, headLen int32, tagLen int32) *TCPClient {
	return &TCPClient{
		conn: conn,
		r:    bufio.NewReader(conn),
		w:    bufio.NewWriter(conn),
		head: headLen,
		tag:  tagLen,
	}
}

func (receiver *TCPClient) LocalAddr() net.Addr {
	return receiver.conn.LocalAddr()
}

func (receiver *TCPClient) RemoteAddr() net.Addr {
	return receiver.conn.RemoteAddr()
}

func (receiver *TCPClient) Close() error {
	return receiver.conn.Close()
}

func (receiver *TCPClient) Write(message []byte, tag int32) (int, error) {
	// 读取消息长度
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.BigEndian, length)
	if err != nil {
		return 0, err
	}
	// 写入消息类型
	err = binary.Write(pkg, binary.BigEndian, tag)
	if err != nil {
		return 0, err
	}
	// 写入消息体
	err = binary.Write(pkg, binary.BigEndian, message)
	if err != nil {
		return 0, err
	}
	write, err := receiver.w.Write(pkg.Bytes())
	if err != nil {
		return 0, err
	}

	err = receiver.w.Flush()
	if err != nil {
		return 0, err
	}
	return write, nil
}

func (receiver *TCPClient) Read() (int32, []byte, error) {
	// Peek返回缓存的一个切片，该切片引用缓存中前N个字节的数据
	lengthByte, err := receiver.r.Peek(int(receiver.head + receiver.tag))
	if err != nil {
		return 0, nil, err
	}
	// 创建buffer缓冲器
	var dataLen int32
	// lengthByte = [head,tag]
	lengthBuff := bytes.NewBuffer(lengthByte[:receiver.head])
	// 将head的内容读到pkgLen中
	err = binary.Read(lengthBuff, binary.BigEndian, &dataLen)
	if err != nil {
		return 0, nil, err
	}
	// 将类型读取到tag中
	var tag int32
	tagBuff := bytes.NewBuffer(lengthByte[receiver.head:])
	err = binary.Read(tagBuff, binary.BigEndian, &tag)
	if err != nil {
		return 0, nil, err
	}
	// 如果当前在缓冲区可以读的内容小于包长度+头部长度+类型长度，则表示这个包是有问题的
	// 这个意思是一个数据包包括[head+tag+data]，即数据包长度=head的长度+tag的长度+data的长度
	// 而data的长度已经从head中读到了，也就是刚才的dataLen
	if int32(receiver.r.Buffered()) < dataLen+receiver.head+receiver.tag {
		return 0, nil, err
	}
	// 读取消息真正的内容
	pack := make([]byte, int(receiver.head+dataLen+receiver.tag))
	// 如果缓存不为空，只能读出缓存中的数据，不会从底层io.Reader中提取数据
	// 如果缓存为空
	// 1、len(p) >= 缓存大小，则跳过缓存，直接从底层io.Reader中读出到p
	// 2、len(p) < 缓存大小，则先将数据从底层io.Reader读到缓存中，再从缓存读到p中
	_, err = receiver.r.Read(pack)
	if err != nil {
		return 0, nil, err
	}
	return tag, pack[receiver.head+receiver.tag:], nil
}
