package net

import (
	"errors"
	"os"
	"time"
)

var errNotImplemented = errors.New("not implemented")

func sysInit() {}

type netFD struct {
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

// lookupProtocol looks up IP protocol name and returns
// the corresponding protocol number.
func lookupProtocol(name string) (proto int, err error) {
	return 0, errNotImplemented
}

func lookupIP(host string) (addrs []IPAddr, err error) {
	return nil, errNotImplemented
}

func lookupHost(host string) (addrs []string, err error) {
	return nil, errNotImplemented
}

func lookupPort(network, service string) (port int, err error) {
	return 0, errNotImplemented
}

func lookupCNAME(name string) (cname string, err error) {
	return "", errNotImplemented
}

func lookupSRV(service, proto, name string) (cname string, addrs []*SRV, err error) {
	return "", nil, errNotImplemented
}

func lookupMX(name string) (mx []*MX, err error) {
	return nil, errNotImplemented
}

func lookupNS(name string) (ns []*NS, err error) {
	return nil, errNotImplemented
}

func lookupTXT(name string) (txt []string, err error) {
	return nil, errNotImplemented
}

func lookupAddr(addr string) (name []string, err error) {
	return nil, errNotImplemented
}

func dial(net string, ra Addr, dialer func(time.Time) (Conn, error), deadline time.Time) (Conn, error) {
	return nil, errNotImplemented
}

type errConn struct {
	conn
}

func (*errConn) ReadFrom(b []byte) (int, Addr, error) {
	return 0, nil, errNotImplemented
}

func (*errConn) WriteTo(b []byte, addr Addr) (int, error) {
	return 0, errNotImplemented
}

type errListener struct{}

func (*errListener) Accept() (Conn, error) {
	return nil, errNotImplemented
}

func (*errListener) Addr() Addr { return nil }

type TCPConn struct {
	errConn
	errListener
}

func dialTCP(net string, laddr, raddr *TCPAddr, deadline time.Time, cancel <-chan struct{}) (*TCPConn, error) {
	return nil, errNotImplemented
}

func ListenTCP(net string, laddr *TCPAddr) (*TCPConn, error) {
	return nil, errNotImplemented
}

type UDPConn struct {
	errConn
}

func dialUDP(net string, laddr, raddr *UDPAddr, deadline time.Time) (*UDPConn, error) {
	return nil, errNotImplemented
}

func ListenUDP(netProto string, laddr *UDPAddr) (*UDPConn, error) {
	return nil, errNotImplemented
}

type IPConn struct {
	errConn
}

func dialIP(netProto string, laddr, raddr *IPAddr, deadline time.Time) (*IPConn, error) {
	return nil, errNotImplemented
}

func ListenIP(netProto string, laddr *IPAddr) (*IPConn, error) {
	return nil, errNotImplemented
}

type UnixConn struct {
	errConn
	errListener
}

func dialUnix(net string, laddr, raddr *UnixAddr, deadline time.Time) (*UnixConn, error) {
	return nil, errNotImplemented
}

func ListenUnix(net string, laddr *UnixAddr) (*UnixConn, error) {
	return nil, errNotImplemented
}

func ListenUnixgram(net string, laddr *UnixAddr) (*UnixConn, error) {
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
