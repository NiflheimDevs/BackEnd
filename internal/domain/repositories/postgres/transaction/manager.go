package transaction

import "context"

type TxManager interface {
	Begin(ctx context.Context) (Tx, error)
}
