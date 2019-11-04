package id

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"math"
	"sync/atomic"
	"time"
)

var base32IDEncoding = base32.NewEncoding("abcdefghijklmnopqrstuv1234567890").WithPadding(base32.NoPadding)

type ID []byte

func (i ID) Hex() string {
	return hex.EncodeToString(i)
}

func (i ID) Base32() string {
	return base32IDEncoding.EncodeToString(i)
}

type IDGenerator struct {
	randomSize int32
	counter    uint32
}

func NewIDGenerator(randomSize int32) (*IDGenerator, error) {
	counter := make([]byte, 4)
	_, err := rand.Read(counter)
	if err != nil {
		return nil, err
	}
	return &IDGenerator{
		randomSize: randomSize,
		counter:    binary.LittleEndian.Uint32(counter),
	}, nil
}

func (i *IDGenerator) New() (ID, error) {
	id := make(ID, i.randomSize+12)
	atomic.AddUint32(&i.counter, 1)
	_, err := rand.Read(id[:i.randomSize])
	if err != nil {
		return nil, err
	}
	c := atomic.AddUint32(&i.counter, 1)
	binary.LittleEndian.PutUint32(id[i.randomSize:i.randomSize+4], c)
	binary.LittleEndian.PutUint64(id[i.randomSize+4:i.randomSize+12], uint64(time.Now().UnixNano()))
	if c >= math.MaxInt32 {
		atomic.StoreUint32(&i.counter, 0)
	}
	return id, nil
}
