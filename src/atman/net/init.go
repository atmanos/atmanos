package net

import (
	"atman/mm"
	"atman/xen"
	"atman/xenstore"
	"fmt"
)

var grantTable = xen.MapGrantTable()

var DefaultDevice *Device

func init() {
	dev, err := initNetworking()
	if err != nil {
		panic(err)
	}

	DefaultDevice = dev
}

type Device struct {
	Backend uint32

	Tx        *xen.FrontendRing
	TxBuffers *BufferPool
	TxGref    xen.Gref

	Rx        *xen.FrontendRing
	RxBuffers *BufferPool
	RxGref    xen.Gref

	EventChannel *xen.EventChannel

	MacAddr []byte
	IPAddr  []byte
}

// initNetworking sets up the default network device.
func initNetworking() (*Device, error) {
	dev := &Device{}

	backendDomID, err := xenstore.Read("device/vif/0/backend-id").ReadUint32()

	if err != nil {
		return nil, fmt.Errorf("atman/net: unable to read device: %s", err)
	}

	dev.Backend = backendDomID
	dev.EventChannel = xen.NewEventChannel(backendDomID)

	txPage := mm.AllocPage()
	dev.Tx = newTxRing(initSharedRing(txPage))
	dev.TxGref = mustGrantAccess(dev.Backend, txPage.Frame, false)
	dev.TxBuffers = NewBufferPool(int(dev.Tx.EntryCount))

	rxPage := mm.AllocPage()
	dev.Rx = newRxRing(initSharedRing(rxPage))
	dev.RxGref = mustGrantAccess(dev.Backend, rxPage.Frame, false)
	dev.RxBuffers = NewBufferPool(int(dev.Rx.EntryCount))

	initRxPages(dev)

	if err := dev.register(); err != nil {
		return nil, fmt.Errorf("atman/net: failed to register device: %s", err)
	}

	if err := dev.finalizeConnection(); err != nil {
		return nil, fmt.Errorf("atman/net: failed to finalize connection: %s", err)
	}

	return dev, nil
}

// initRxPages allocates buffers for receiving rx packets
// and sends them to the backend.
func initRxPages(dev *Device) {
	for {
		buf, ok := dev.RxBuffers.Get()
		if !ok {
			break
		}

		buf.Gref = mustGrantAccess(dev.Backend, buf.Page.Frame, false)

		req := (*NetifRxRequest)(dev.Rx.NextRequest())
		req.ID = uint16(buf.ID)
		req.Gref = buf.Gref
	}

	if notify := dev.Rx.PushRequests(); notify {
		dev.EventChannel.Notify()
	}
}

// register registers the device in the Xen Store.
func (dev *Device) register() error {
	for committed := false; !committed; {
		tx, err := xenstore.TransactionStart()
		if err != nil {
			return err
		}

		tx.WriteInt(dev.xenstorePath("tx-ring-ref"), int(dev.TxGref))
		tx.WriteInt(dev.xenstorePath("rx-ring-ref"), int(dev.RxGref))
		tx.WriteInt(dev.xenstorePath("event-channel"), int(dev.EventChannel.Port))
		tx.WriteInt(dev.xenstorePath("request-rx-copy"), 1)
		tx.WriteInt(dev.xenstorePath("state"), xenstore.StateConnected)

		committed, err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

// finalizeConnection waits for the backend connection to be ready,
// and reads the device's intended mac and ip addresses.
func (dev *Device) finalizeConnection() error {
	backend, err := xenstore.Read(dev.xenstorePath("backend")).ReadBytes()
	if err != nil {
		return fmt.Errorf("atman/net: failed to read backend path: %s", err)
	}

	state, err := xenstore.Read(string(backend) + "/state").ReadUint32()
	if err != nil {
		return fmt.Errorf("atman/net: failed to read backend state: %s", err)
	}

	if state != xenstore.StateConnected {
		return fmt.Errorf("atman/net: backend not connected (state=%v)", state)
	}

	ip, err := xenstore.Read(string(backend) + "/ip").ReadBytes()
	if err == nil {
		dev.IPAddr = ip
	}

	mac, err := xenstore.Read(dev.xenstorePath("mac")).ReadBytes()
	if err != nil {
		return fmt.Errorf("atman/net: failed to read mac: %s", err)
	}
	dev.MacAddr = mac

	return nil
}

func (dev *Device) xenstorePath(path string) string {
	return "device/vif/0/" + path
}

func (dev *Device) SendTxBuffer(buf *Buffer, size int) {
	buf.Gref = mustGrantAccess(dev.Backend, buf.Page.Frame, true)

	req := (*NetifTxRequest)(dev.Tx.NextRequest())
	req.Gref = buf.Gref
	req.Offset = 0
	req.Flags = 0
	req.ID = uint16(buf.ID)
	req.Size = uint16(size)
}

func mustGrantAccess(dom uint32, frame uintptr, readonly bool) xen.Gref {
	gref, ok := grantTable.GrantAccess(uint16(dom), frame, readonly)

	if !ok {
		panic("unable to grant access to page")
	}

	return gref
}
