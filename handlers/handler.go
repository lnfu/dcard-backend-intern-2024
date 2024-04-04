package handlers

import (
	"fmt"

	"github.com/lnfu/dcard-intern/cache"
	db "github.com/lnfu/dcard-intern/db/sqlc"
)

type Handler struct {
	databaseQueries *db.Queries
	cac             *cache.Cache
}

func NewHandler(db *db.Queries, cac *cache.Cache) *Handler {
	return &Handler{db, cac}
}

type InvalidQueryParameterError struct {
	ParameterName string
	Reason        string
}

func (e InvalidQueryParameterError) Error() string {
	return fmt.Sprintf("Invalid query parameter '%s': %s", e.ParameterName, e.Reason)
}
