package handlers

import (
	"fmt"

	db "github.com/lnfu/dcard-intern/db/sqlc"
)

type Handler struct {
	databaseQueries *db.Queries
}

func NewHandler(db *db.Queries) *Handler {
	return &Handler{db}
}

type InvalidQueryParameterError struct {
	ParameterName string
	Reason        string
}

func (e InvalidQueryParameterError) Error() string {
	return fmt.Sprintf("Invalid query parameter '%s': %s", e.ParameterName, e.Reason)
}
