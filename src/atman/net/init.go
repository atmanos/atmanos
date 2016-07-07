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
	DefaultDevice = initNetworking()
}

type buffer struct {
	Gref xen.Gref
	*mm.Page
}

type Device struct {
	Backend uint32

	Tx     *xen.FrontendRing
	TxGref xen.Gref

	Rx        *xen.FrontendRing
	RxBuffers []buffer
	RxGref    xen.Gref

	EventChannel *xen.EventChannel

	MacAddr []byte
	IPAddr  []byte
}

func initNetworking() *Device {
	dev := &Device{}

	backendDomID, err := xenstore.Read("device/vif/0/backend-id").ReadUint32()

	if err != nil {
		println("Unable to read device/vif/0/backend-id")
		panic(err)
	}

	dev.Backend = backendDomID
	dev.EventChannel = xen.NewEventChannel(backendDomID)

	// setup tx freelist
	txPage := mm.AllocPage()
	dev.Tx = newTxRing(initSharedRing(txPage))
	dev.TxGref = mustGrantAccess(dev.Backend, txPage.Frame, false)

	rxPage := mm.AllocPage()
	dev.Rx = newRxRing(initSharedRing(rxPage))
	dev.RxGref = mustGrantAccess(dev.Backend, rxPage.Frame, false)
	dev.RxBuffers = make([]buffer, dev.Rx.EntryCount)

	initRxPages(dev)

	if err := dev.register(); err != nil {
		println("Failed to register device: ", err.Error())
		return nil
	}

	dev.MacAddr, _ = xenstore.Read(dev.xenstorePath("mac")).ReadBytes()

	backend, _ := xenstore.Read(dev.xenstorePath("backend")).ReadBytes()

	state, _ := xenstore.Read(string(backend) + "/state").ReadUint32()
	if state != xenstore.StateConnected {
		fmt.Println("net: backend state=%v waiting for connection", state)
		return nil
	}

	ip, err := xenstore.Read(string(backend) + "/ip").ReadBytes()
	if err == nil {
		dev.IPAddr = ip
	}

	return dev
}

func mustGrantAccess(dom uint32, frame uintptr, readonly bool) xen.Gref {
	gref, ok := grantTable.GrantAccess(uint16(dom), frame, readonly)

	if !ok {
		panic("unable to grant access to page")
	}

	return gref
}

func initRxPages(dev *Device) {
	for i := range dev.RxBuffers {
		buf := &dev.RxBuffers[i]
		buf.Page = mm.AllocPage()
		buf.Gref = mustGrantAccess(dev.Backend, buf.Page.Frame, false)

		req := (*NetifRxRequest)(dev.Rx.NextRequest())
		req.ID = uint16(i)
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

func (dev *Device) xenstorePath(path string) string {
	return "device/vif/0/" + path
}
