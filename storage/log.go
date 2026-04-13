package storage

import (
	"encoding/binary"
	"os"
)

type Log struct {
	file   *os.File
	offset int64
}

func NewLog(path string) (*Log, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	info, _ := file.Stat()

	return &Log{
		file:   file,
		offset: info.Size(),
	}, nil
}

func (l *Log) Append(message []byte) (int64, error) {
	currentOffset := l.offset

	// Write message length (8 bytes)
	err := binary.Write(l.file, binary.BigEndian, int64(len(message)))
	if err != nil {
		return 0, err
	}

	// Write message
	_, err = l.file.Write(message)
	if err != nil {
		return 0, err
	}

	l.offset += 8 + int64(len(message))

	return currentOffset, nil
}

func (l *Log) Read(offset int64) ([]byte, int64, error) {
	_, err := l.file.Seek(offset, 0)
	if err != nil {
		return nil, 0, err
	}

	var msgLen int64
	err = binary.Read(l.file, binary.BigEndian, &msgLen)
	if err != nil {
		return nil, 0, err
	}

	msg := make([]byte, msgLen)
	_, err = l.file.Read(msg)
	if err != nil {
		return nil, 0, err
	}

	nextOffset := offset + 8 + msgLen
	return msg, nextOffset, nil
}
