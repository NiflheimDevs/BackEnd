package repositoriesimpl

import "github.com/jackc/pgx/v5/pgxpool"

type TeamRepo struct {
    PG *pgxpool.Pool
}

func NewteamRepo (
    PG *pgxpool.Pool,
    ) *TeamRepo {
    return &TeamRepo{
        PG: PG,
    }
}
