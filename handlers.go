package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/lnfu/dcard-intern/db/sqlc"
)

func NullStringFromString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func NullInt32FromString(s string) (sql.NullInt32, error) {
	if s == "" {
		return sql.NullInt32{Valid: false}, nil
	}
	temp, err := strconv.Atoi(s) // TODO error handling (age 輸入給錯)
	if err != nil {
		return sql.NullInt32{Valid: false}, err
	}
	return sql.NullInt32{Int32: int32(temp), Valid: true}, nil
}

func (app *application) getAdvertisementFilters(ctx *gin.Context) (sql.NullInt32, sql.NullString, sql.NullString, sql.NullString, int, int) {
	// TODO 可能要做多選參數 (e.g., gender=M,F)
	var (
		age                       sql.NullInt32
		gender, country, platform sql.NullString
		offset, limit             int
	)

	age, _ = NullInt32FromString(ctx.Query("age")) // TODO error handling (age 輸入給錯)
	gender = NullStringFromString(ctx.Query("gender"))
	country = NullStringFromString(ctx.Query("country"))
	platform = NullStringFromString(ctx.Query("platform"))

	offset, _ = strconv.Atoi(ctx.Query("offset")) // TODO error handling (offset 輸入給錯)
	limit, _ = strconv.Atoi(ctx.Query("limit"))   // TODO error handling (limit 輸入給錯)
	if limit == 0 {
		limit = 5
	}

	return age, gender, country, platform, offset, limit
}
func (app *application) getAdvertisementHandler(ctx *gin.Context) {
	age, gender, country, platform, offset, limit := app.getAdvertisementFilters(ctx)

	ads, err := app.databaseQueries.GetAdvertisements(ctx, db.GetAdvertisementsParams{
		Age:      age,
		Gender:   gender,
		Country:  country,
		Platform: platform,
		Offset:   int32(offset),
		Limit:    int32(limit),
	})

	if err != nil {
		app.errorLogger.Print(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": ads,
	})
}

type CreateAdvertisementRequest struct {
	Title      string
	StartAt    time.Time
	EndAt      time.Time
	Conditions []AdvertisementCondition
}

type AdvertisementCondition struct {
	AgeStart sql.NullInt32
	AgeEnd   sql.NullInt32
	Gender   sql.NullString
	Country  sql.NullString
	Platform sql.NullString
}

func (app *application) createAdvertisementHandler(ctx *gin.Context) {
	body := CreateAdvertisementRequest{}
	ctx.BindJSON(&body)

	advertisementId, _ := app.databaseQueries.CreateAdvertisement(ctx, db.CreateAdvertisementParams{
		Title:   body.Title,
		StartAt: body.StartAt,
		EndAt:   body.EndAt,
	})

	for _, condition := range body.Conditions {
		conditionId, _ := app.databaseQueries.CreateCondition(ctx, db.CreateConditionParams{
			AgeStart: condition.AgeStart,
			AgeEnd:   condition.AgeEnd,
		})
		if condition.Gender.Valid {
			app.databaseQueries.CreateConditionGender(ctx, db.CreateConditionGenderParams{
				ConditionID: int32(conditionId),
				Gender:      condition.Gender.String,
			})
		}
		if condition.Country.Valid {
			app.databaseQueries.CreateConditionCountry(ctx, db.CreateConditionCountryParams{
				ConditionID: int32(conditionId),
				Country:     condition.Country.String,
			})
		}
		if condition.Platform.Valid {
			app.databaseQueries.CreateConditionPlatform(ctx, db.CreateConditionPlatformParams{
				ConditionID: int32(conditionId),
				Platform:    condition.Platform.String,
			})
		}
		app.databaseQueries.CreateAdvertisementCondition(ctx, db.CreateAdvertisementConditionParams{
			AdvertisementID: int32(advertisementId),
			ConditionID:     int32(conditionId),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
