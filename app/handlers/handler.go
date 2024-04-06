package handlers

import (
	"context"
	"fmt"
	"log"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/lnfu/dcard-intern/app/cache"
	db "github.com/lnfu/dcard-intern/app/db/sqlc"
)

var ctx = context.Background()

type Handler struct {
	databaseQueries *db.Queries
	cac             *cache.Cache
	genderSet       mapset.Set[string]
	countrySet      mapset.Set[string]
	platformSet     mapset.Set[string]
}

func NewHandler(db *db.Queries, cac *cache.Cache) *Handler {
	genders, err := db.GetAllGenders(ctx)
	if err != nil {
		log.Fatalln("Database error", err.Error())
	}
	genderSet := mapset.NewSet[string]()
	for _, gender := range genders {
		genderSet.Add(gender)
	}

	countries, err := db.GetAllCountries(ctx)
	if err != nil {
		log.Fatalln("Database error", err.Error())
	}
	countrySet := mapset.NewSet[string]()
	for _, country := range countries {
		countrySet.Add(country)
	}

	platforms, err := db.GetAllPlatforms(ctx)
	if err != nil {
		log.Fatalln("Database error", err.Error())
	}
	platformSet := mapset.NewSet[string]()
	for _, platform := range platforms {
		platformSet.Add(platform)
	}

	return &Handler{db, cac, genderSet, countrySet, platformSet}
}

type InvalidQueryParameterError struct {
	ParameterName string
	Reason        string
}

func (e InvalidQueryParameterError) Error() string {
	return fmt.Sprintf("Invalid query parameter '%s': %s", e.ParameterName, e.Reason)
}
