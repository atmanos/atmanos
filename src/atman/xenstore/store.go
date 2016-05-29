package xenstore

import (
	"atman/ring"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
	"unsafe"
)

type _store struct {
	port uint32

	req *ring.Buffer
	rsp *ring.Buffer
}

var store *_store

func init() {
	type xenstoreInterface struct {
		req, rsp [1024]byte

		reqCons, reqProd uint32
		rspCons, rspProd uint32

		serverFeatures uint32
		connState      uint32
	}

	r, port := runtime.LoadXenStore()

	xenRing := (*xenstoreInterface)(r)

	store = &_store{
		port: port,
		req: &ring.Buffer{
			Data:      xenRing.req[:],
			ReaderPos: &xenRing.reqCons,
			WriterPos: &xenRing.reqProd,
		},
		rsp: &ring.Buffer{
			Data:      xenRing.rsp[:],
			ReaderPos: &xenRing.rspCons,
			WriterPos: &xenRing.rspProd,
		},
	}

	runtime.BindEventHandler(port)
	go loop()
}

func loop() {
	for {
		runtime.WaitEvent(store.port)
		handleEvent()
	}
}

func handleEvent() {
	var data = make([]byte, MessageHeaderSize)

	if n := store.rsp.Read(data); n != len(data) {
		fmt.Printf("rsp read n=%v data=%v\n", n, data)
		return
	}

	var msg MessageHeader

	if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &msg); err != nil {
		println(err)
		return
	}

	fmt.Printf("%#v\n", msg)
	data = make([]byte, msg.Length)
	store.rsp.Read(data)
	fmt.Printf("%s", data[:len(data)-1])

	for {
		n := store.rsp.Read(data)
		if n == 0 {
			// consume rest
			return
		}
	}
}

const MessageHeaderSize = int(unsafe.Sizeof(MessageHeader{}))

type MessageHeader struct {
	Type        uint32
	RequestID   uint32
	Transaction uint32
	Length      uint32
}

const (
	TypeDebug = iota
	TypeDirectory
	TypeRead
	TypeGetPerms
	TypeWatch
	TypeUnwatch
	TypeTransactionStart
	TypeTransactionEnd
	TypeIntroduce
	TypeRelease
	TypeGetDomainPath
	TypeWrite
	TypeMkdir
	TypeRm
	TypeSetPerms
	TypeWatchEvent
	TypeError
	TypeIsDomainIntroduced
	TypeResume
	TypeSetTarget
	TypeRestrict
	TypeResetWatches

	TypeInvalid = 0xffff
)

const (
	WatchPath = iota
	WatchToken
)

func Read(path string) ([]byte, error) {
	buf := &bytes.Buffer{}
	msg := MessageHeader{
		Type:        TypeDirectory,
		RequestID:   1,
		Transaction: 0,
		Length:      uint32(len(path) + 1),
	}

	if err := binary.Write(buf, binary.LittleEndian, msg); err != nil {
		return nil, err
	}

	buf.WriteString(path)
	buf.WriteByte(0)

	if n := store.req.Write(buf.Bytes()); n != buf.Len() {
		return nil, io.ErrShortWrite
	}

	runtime.NotifyEventChannel(store.port)

	return nil, nil
}

func Debug(s string) error {
	buf := &bytes.Buffer{}
	msg := MessageHeader{
		Type:        TypeDebug,
		RequestID:   1,
		Transaction: 0,
		Length:      uint32(len("print") + 1 + len(s) + 1),
	}

	if err := binary.Write(buf, binary.LittleEndian, msg); err != nil {
		return err
	}

	buf.WriteString("print")
	buf.WriteByte(0)
	buf.WriteString(s)
	buf.WriteByte(0)

	n := store.req.Write(buf.Bytes())
	fmt.Printf("Write n=%v data=%#v\n", buf.Bytes())

	if n != buf.Len() {
		return io.ErrShortWrite
	}

	runtime.NotifyEventChannel(store.port)
	return nil
}
