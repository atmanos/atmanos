package net

import (
	"atman/mm"
	"atman/xen"
	"atman/xenstore"
	"fmt"
)

var grantTable = xen.MapGrantTable()

func init() {
	initNetworking()
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
}

func initNetworking() {
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
	dev.Tx = newTxRing(txPage)
	dev.TxGref = mustGrantAccess(dev.Backend, txPage.Frame, false)

	rxPage := mm.AllocPage()
	dev.Rx = newRxRing(rxPage)
	dev.RxGref = mustGrantAccess(dev.Backend, rxPage.Frame, false)
	dev.RxBuffers = make([]buffer, dev.Rx.EntryCount)
	initRxPages(dev)

	if err := dev.register(); err != nil {
		println("Failed to register device: ", err.Error())
		return
	}

	backend, _ := xenstore.Read(dev.xenstorePath("backend")).ReadBytes()
	mac, _ := xenstore.Read(dev.xenstorePath("mac")).ReadBytes()

	fmt.Printf("net: backend=%q mac=%v (%q)\n", backend, mac, mac)

	state, _ := xenstore.Read(string(backend) + "/state").ReadUint32()
	if state != xenstore.StateConnected {
		fmt.Println("net: backend state=%v waiting for connection", state)
		return
	}

	ip, _ := xenstore.Read(string(backend) + "/ip").ReadBytes()
	fmt.Printf("net: ip=%v (%q)\n", ip, ip)
}

func mustGrantAccess(dom uint32, frame uintptr, readonly bool) xen.Gref {
	gref, ok := grantTable.GrantAccess(uint16(dom), frame, readonly)

	if !ok {
		panic("unable to grant access to page")
	}

	return gref
}

func initRxPages(dev *Device) {
	for i, buf := range dev.RxBuffers {
		buf.Page = mm.AllocPage()
		buf.Gref = mustGrantAccess(dev.Backend, buf.Page.Frame, false)

		req := (*netifRxRequest)(dev.Rx.NextRequest())
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
