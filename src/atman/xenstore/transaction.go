package xenstore

type Transaction struct {
	id  uint32
	err error
}

func TransactionStart() (*Transaction, error) {
	tx := &Transaction{}

	if err := tx.start(); err != nil {
		return nil, err
	}

	return tx, nil
}

func (tx *Transaction) start() error {
	req := NewRequest(TypeTransactionStart, 0)
	req.WriteString("") // xenstored requires an empty argument

	rsp := Send(req)
	if err := rsp.Err(); err != nil {
		return err
	}

	id, err := rsp.ReadUint32()
	if err != nil {
		return err
	}

	tx.id = id
	return nil
}

func (tx *Transaction) WriteInt(path string, i int) {
	if tx.err != nil {
		return
	}

	req := NewRequest(TypeWrite, tx.id)
	req.WriteUint32(uint32(i))

	rsp := Send(req)
	if err := rsp.Err(); err != nil {
		tx.abortWith(err)
		return
	}
}

func (tx *Transaction) SwitchState(path string, state int) {
}

func (tx *Transaction) abortWith(err error) {
	tx.err = err

	req := NewRequest(TypeTransactionEnd, tx.id)
	req.WriteString("F")

	Send(req)
}

func (tx *Transaction) Commit() (committed bool, err error) {
	if tx.err != nil {
		return false, tx.err
	}

	req := NewRequest(TypeTransactionEnd, tx.id)
	req.WriteString("T")

	err = Send(req).Err()
	return err == ErrRetry, err
}
