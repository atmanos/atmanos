package xenstore

import (
	"strconv"
)

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

	s, err := rsp.ReadString()
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}

	tx.id = uint32(id)
	return nil
}

func (tx *Transaction) ReadInt(path string) (int, error) {
	if tx.err != nil {
		return 0, tx.err
	}

	req := NewRequest(TypeRead, tx.id)
	req.WriteString(path)

	rsp := Send(req)
	if err := rsp.Err(); err != nil {
		tx.abortWith(annotateError(TypeRead, path, err))
		return 0, err
	}

	i, err := rsp.ReadUint32()
	return int(i), err
}

func (tx *Transaction) WriteInt(path string, i int) {
	if tx.err != nil {
		return
	}

	req := NewRequest(TypeWrite, tx.id)
	req.WriteString(path)
	req.WriteUint32(uint32(i))

	rsp := Send(req)
	if err := rsp.Err(); err != nil {
		tx.abortWith(annotateError(TypeWrite, path, err))
		return
	}
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

	if isRetry(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
