package tests

import "context"

type TxMock struct {
}

func (t TxMock) Commit(_ context.Context) error {
	return nil
}

func (t TxMock) Rollback(_ context.Context) error {
	return nil
}
