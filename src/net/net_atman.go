package net

import (
	"context"
	"errors"
	"io"
	"os"
	"time"
)

var errNotImplemented = errors.New("not implemented")

func sysInit() {}

type netFD struct {
	fdmu fdMutex

	net          string
	laddr, raddr Addr
}

func (*netFD) Close() error {
	return errNotImplemented
}

func (*netFD) Read(b []byte) (n int, err error) {
	return 0, errNotImplemented
}

func (*netFD) Write(b []byte) (n int, err error) {
	return 0, errNotImplemented
}

func (*netFD) destroy() {}

func (*netFD) dup() (*os.File, error) {
	return nil, errNotImplemented
}

func (fd *netFD) setDeadline(t time.Time) error {
	return errNotImplemented
}

func (fd *netFD) setReadDeadline(t time.Time) error {
	return errNotImplemented
}

func (fd *netFD) setWriteDeadline(t time.Time) error {
	return errNotImplemented
}

func (*netFD) closeRead() error {
	return errNotImplemented
}

func (*netFD) closeWrite() error {
	return errNotImplemented
}

func setReadBuffer(fd *netFD, bytes int) error {
	return errNotImplemented
}

func setWriteBuffer(fd *netFD, bytes int) error {
	return errNotImplemented
}

func setKeepAlive(fd *netFD, keepalive bool) error {
	return errNotImplemented
}

func setKeepAlivePeriod(fd *netFD, d time.Duration) error {
	return errNotImplemented
}

func setLinger(fd *netFD, sec int) error {
	return errNotImplemented
}

func setNoDelay(fd *netFD, noDelay bool) error {
	return errNotImplemented
}

// lookupProtocol looks up IP protocol name and returns
// the corresponding protocol number.
func lookupProtocol(_ context.Context, name string) (proto int, err error) {
	return 0, errNotImplemented
}

func (*Resolver) lookupIP(_ context.Context, host string) (addrs []IPAddr, err error) {
	return nil, errNotImplemented
}

func (*Resolver) lookupHost(_ context.Context, host string) (addrs []string, err error) {
	return nil, errNotImplemented
}

func (*Resolver) lookupPort(_ context.Context, network, service string) (port int, err error) {
	return 0, errNotImplemented
}

func (*Resolver) lookupCNAME(_ context.Context, name string) (cname string, err error) {
	return "", errNotImplemented
}

func (*Resolver) lookupSRV(_ context.Context, service, proto, name string) (cname string, addrs []*SRV, err error) {
	return "", nil, errNotImplemented
}

func (*Resolver) lookupMX(_ context.Context, name string) (mx []*MX, err error) {
	return nil, errNotImplemented
}

func (*Resolver) lookupNS(_ context.Context, name string) (ns []*NS, err error) {
	return nil, errNotImplemented
}

func (*Resolver) lookupTXT(_ context.Context, name string) (txt []string, err error) {
	return nil, errNotImplemented
}

func (*Resolver) lookupAddr(_ context.Context, addr string) (name []string, err error) {
	return nil, errNotImplemented
}

func dial(net string, ra Addr, dialer func(time.Time) (Conn, error), deadline time.Time) (Conn, error) {
	return nil, errNotImplemented
}

func dialTCP(_ context.Context, net string, laddr, raddr *TCPAddr) (*TCPConn, error) {
	return nil, errNotImplemented
}

func listenTCP(_ context.Context, net string, laddr *TCPAddr) (*TCPListener, error) {
	return nil, errNotImplemented
}

func dialUDP(_ context.Context, net string, laddr, raddr *UDPAddr) (*UDPConn, error) {
	return nil, errNotImplemented
}

func listenUDP(_ context.Context, netProto string, laddr *UDPAddr) (*UDPConn, error) {
	return nil, errNotImplemented
}

func listenMulticastUDP(_ context.Context, network string, ifi *Interface, gaddr *UDPAddr) (*UDPConn, error) {
	return nil, errNotImplemented
}

func dialIP(_ context.Context, netProto string, laddr, raddr *IPAddr) (*IPConn, error) {
	return nil, errNotImplemented
}

func listenIP(_ context.Context, netProto string, laddr *IPAddr) (*IPConn, error) {
	return nil, errNotImplemented
}

func dialUnix(_ context.Context, net string, laddr, raddr *UnixAddr) (*UnixConn, error) {
	return nil, errNotImplemented
}

func listenUnix(_ context.Context, net string, laddr *UnixAddr) (*UnixListener, error) {
	return nil, errNotImplemented
}

func listenUnixgram(_ context.Context, net string, laddr *UnixAddr) (*UnixConn, error) {
	return nil, errNotImplemented
}

func fileConn(f *os.File) (Conn, error) {
	return nil, errNotImplemented
}

func filePacketConn(f *os.File) (PacketConn, error) {
	return nil, errNotImplemented
}

func fileListener(f *os.File) (Listener, error) {
	return nil, errNotImplemented
}

func interfaceTable(ifindex int) ([]Interface, error) {
	return nil, errNotImplemented
}

func interfaceAddrTable(ifi *Interface) ([]Addr, error) {
	return nil, errNotImplemented
}

func interfaceMulticastAddrTable(ifi *Interface) ([]Addr, error) {
	return nil, errNotImplemented
}

func probeIPv4Stack() bool {
	return false
}

func probeIPv6Stack() (supportsIPv6, supportsIPv4map bool) {
	return false, false
}

func maxListenerBacklog() int {
	return -1
}

func (*IPConn) writeTo(b []byte, addr *IPAddr) (int, error) {
	return 0, errNotImplemented
}

func (*IPConn) readFrom(b []byte) (int, *IPAddr, error) {
	return 0, nil, errNotImplemented
}

func (*IPConn) readMsg(b, oob []byte) (n, oobn, flags int, addr *IPAddr, err error) {
	return 0, 0, 0, nil, errNotImplemented
}

func (*IPConn) writeMsg(b, oob []byte, addr *IPAddr) (n, oobn int, err error) {
	return 0, 0, errNotImplemented
}

func (*TCPConn) readFrom(r io.Reader) (int64, error) { return 0, errNotImplemented }

func (*TCPListener) ok() bool                  { return false }
func (*TCPListener) accept() (*TCPConn, error) { return nil, errNotImplemented }
func (*TCPListener) close() error              { return errNotImplemented }
func (*TCPListener) file() (*os.File, error)   { return nil, errNotImplemented }

func (*UDPConn) readFrom(b []byte) (n int, addr *UDPAddr, err error) {
	return 0, nil, errNotImplemented
}

func (*UDPConn) readMsg(b, oob []byte) (n, oobn, flags int, addr *UDPAddr, err error) {
	return 0, 0, 0, nil, errNotImplemented
}

func (*UDPConn) writeTo(b []byte, addr *UDPAddr) (int, error) {
	return 0, errNotImplemented
}

func (*UDPConn) writeMsg(b, oob []byte, addr *UDPAddr) (n, oobn int, err error) {
	return 0, 0, errNotImplemented
}

func (*UnixListener) accept() (*UnixConn, error) { return nil, errNotImplemented }
func (*UnixListener) close() error               { return errNotImplemented }
func (*UnixListener) file() (*os.File, error)    { return nil, errNotImplemented }

func (*UnixConn) readFrom(b []byte) (n int, addr *UnixAddr, err error) {
	return 0, nil, errNotImplemented
}

func (*UnixConn) readMsg(b, oob []byte) (n, oobn, flags int, addr *UnixAddr, err error) {
	return 0, 0, 0, nil, errNotImplemented
}

func (*UnixConn) writeTo(b []byte, addr *UnixAddr) (int, error) {
	return 0, errNotImplemented
}

func (*UnixConn) writeMsg(b, oob []byte, addr *UnixAddr) (n, oobn int, err error) {
	return 0, 0, errNotImplemented
}
