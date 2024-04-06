package handlers

import (
	"context"
	"fmt"
	"log"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/lnfu/dcard-intern/app/cache"
	sqlc "github.com/lnfu/dcard-intern/app/models/sqlc"
)

func Int32Ptr(i int32) *int32    { return &i }
func StringPtr(s string) *string { return &s }

var ctx = context.Background()

type Handler struct {
	databaseQueries *sqlc.Queries
	cac             *cache.Cache
	genderSet       mapset.Set[string]
	countrySet      mapset.Set[string]
	platformSet     mapset.Set[string]
}

func NewHandler(db *sqlc.Queries, cac *cache.Cache) *Handler {
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
