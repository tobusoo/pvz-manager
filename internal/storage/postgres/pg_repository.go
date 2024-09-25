package postgres

type PgRepository struct {
	txManager TransactionManager
}

func NewRepoPG(tx TransactionManager) *PgRepository {
	return &PgRepository{
		txManager: tx,
	}
}
