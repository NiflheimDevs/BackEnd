package seed

import (
	"context"
	"log"

	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type Seeder struct {
	Manager transaction.TxManager
}

func NewSeeder(manager transaction.TxManager) *Seeder {
	return &Seeder{
		Manager: manager,
	}
}

func (s *Seeder) SeedDaddy() {
	tx, err := s.Manager.Begin(context.Background())
	if err != nil {
		panic(err)
	}

	defer func() {
		if rec := recover(); rec != nil {
			_ = tx.Rollback(context.Background())
			log.Println("wow")
			panic(rec)
		}
	}()

	SeedPerms(tx)
	log.Println("permission seeded")
	SeedRoles(tx)
	log.Println("role seeded")

	if err = tx.Commit(context.Background()); err != nil {
		panic(err)
	}

}
