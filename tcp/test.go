package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/prometheus/common/log"
	"net"
	"testing"
	"time"
)

func Test(t *testing.T) {
	conn, err := net.Dial("tcp", "")
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()
	go func() {
		data, err := Encode("2")
		if err == nil {
			time.Sleep(time.Second * 4)
			_, err := conn.Write(data)
			if err != nil {
				log.Error(err)
			}
		}
	}()
	reader := bufio.NewReader(conn)
	for {
		tag, data, err := Read(reader, 4, 4)
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Println(tag)
		fmt.Println(string(data))
	}
}

func Encode(message string) ([]byte, error) {
	// 读取消息长度
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头（消息头一般是包长）
	err := binary.Write(pkg, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息类型
	err = binary.Write(pkg, binary.BigEndian, int32(0x1))
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.BigEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func Read(r *bufio.Reader, headLen, tagLen int32) (int32, []byte, error) {
	lengthByte, err := r.Peek(int(headLen + tagLen))
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	var length int32
	lengthBuff := bytes.NewBuffer(lengthByte[:headLen])
	err = binary.Read(lengthBuff, binary.BigEndian, &length)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	var tag int32
	tagBuff := bytes.NewBuffer(lengthByte[headLen:])
	err = binary.Read(tagBuff, binary.BigEndian, &tag)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	if int32(r.Buffered()) < length+headLen+tagLen {
		return 0, nil, nil
	}
	pack := make([]byte, int(headLen+length+tagLen))
	_, err = r.Read(pack)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	return tag, pack[headLen+tagLen:], nil
}
