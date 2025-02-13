package repository

import (
	"github.com/jackc/pgx/v5"
)

type IAnomalyRepository interface {
}

type PostgreSQLAnomalyRepository struct {
	conn *pgx.Conn
}

func NewPostgreSQLAnomalyRepository(conn *pgx.Conn) IAnomalyRepository {
	return &PostgreSQLAnomalyRepository{conn: conn}
}
