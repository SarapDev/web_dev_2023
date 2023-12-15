package api

import (
	"database/sql"
	"github.kolesa-team.org/backend/go-module/chi"
)

func HealthHandlerFunc(db *sql.DB) chi.HealthHandlerFunc {
	return db.Ping
}
