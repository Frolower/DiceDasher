package repository

import (
	"backend/internal/auth"
	"backend/internal/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ pool *pgxpool.Pool }

var _ user.UserStore = (*Repository)(nil)

var _ auth.UserStore = (*Repository)(nil)
var _ auth.SessionStore = (*Repository)(nil)

func New(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }
