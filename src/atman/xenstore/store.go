package xenstore

import (
	"atman/ring"
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"
)

type xenStore struct {
	port uint32

	sync.Mutex

	req io.Writer
	rsp io.Reader

	requests map[uint32]chan *Response
}

var store *xenStore

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

	store = &xenStore{
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
		requests: make(map[uint32]chan *Response),
	}

	runtime.BindEventHandler(port)
	go loop()
}

func loop() {
	for {
		runtime.WaitEvent(store.port)
		handleEvents()
	}
}

func handleEvents() {
	for {
		var rsp Response

		if err := rsp.readFrom(store.rsp); err != nil {
			return
		}

		store.handleResponse(&rsp)
	}
}

const MessageHeaderSize = unsafe.Sizeof(MessageHeader{})

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

type Request struct {
	Type        uint32
	ID          uint32
	Transaction uint32

	buf *bytes.Buffer
}

func NewRequest(reqType uint32, tx uint32) *Request {
	req := &Request{
		Type:        reqType,
		ID:          nextRequestID(),
		Transaction: tx,

		buf: &bytes.Buffer{},
	}

	return req
}

func (r *Request) header() MessageHeader {
	return MessageHeader{
		Type:        r.Type,
		RequestID:   r.ID,
		Transaction: r.Transaction,
		Length:      uint32(r.buf.Len()),
	}
}

func (r *Request) WriteString(s string) {
	r.buf.WriteString(s)
	r.buf.WriteByte(0)
}

func (r *Request) WriteBytes(b []byte) {
	r.buf.Write(b)
}

func (r *Request) WriteUint32(i uint32) {
	r.buf.WriteString(strconv.Itoa(int(i)))
}

type Response struct {
	Type      uint32
	RequestID uint32

	buf *bytes.Buffer
}

func (rsp *Response) ReadString() (string, error) {
	return rsp.buf.ReadString(0)
}

func (rsp *Response) ReadUint32() (uint32, error) {
	b, err := rsp.ReadBytes()
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(string(b), 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(i), nil
}

func (rsp *Response) ReadBytes() ([]byte, error) {
	return ioutil.ReadAll(rsp.buf)
}

func (rsp *Response) readFrom(r io.Reader) error {
	var header MessageHeader

	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return err
	}

	rsp.Type = header.Type
	rsp.RequestID = header.RequestID

	data := make([]byte, header.Length)
	io.ReadFull(r, data)

	rsp.buf = bytes.NewBuffer(data)
	return nil
}

func Send(req *Request) *Response {
	done := make(chan *Response, 1)

	store.sendRequest(done, req)

	return <-done
}

func (s *xenStore) sendRequest(done chan *Response, req *Request) {
	s.Lock()
	defer s.Unlock()

	s.requests[req.ID] = done

	header := req.header()

	binary.Write(s.req, binary.LittleEndian, header)
	req.buf.WriteTo(s.req)

	runtime.NotifyEventChannel(s.port)
}

func (s *xenStore) handleResponse(rsp *Response) {
	store.Lock()
	defer store.Unlock()

	c := store.requests[rsp.RequestID]
	delete(store.requests, rsp.RequestID)

	c <- rsp
}

var requestID uint32

func nextRequestID() uint32 {
	return atomic.AddUint32(&requestID, 1)
}
