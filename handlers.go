package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/lnfu/dcard-intern/db/sqlc"

	"github.com/lnfu/dcard-intern/utils"
)

type CreateAdvertisementForm struct {
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

type InvalidQueryParameterError struct {
	ParameterName string
	Reason        string
}

func (e InvalidQueryParameterError) Error() string {
	return fmt.Sprintf("Invalid query parameter '%s': %s", e.ParameterName, e.Reason)
}

func NewInvalidQueryParameterError(parameterName, reason string) InvalidQueryParameterError {
	return InvalidQueryParameterError{
		ParameterName: parameterName,
		Reason:        reason,
	}
}

func (app *application) getAdvertisementQueryParameters(ctx *gin.Context) (sql.NullInt32, sql.NullString, sql.NullString, sql.NullString, int, int, error) {
	// TODO 可能要做多選參數 (e.g., gender=M,F)
	var (
		age                       sql.NullInt32
		gender, country, platform sql.NullString
		offset                    int = 0
		limit                     int = 5
		err                       error
	)

	age, err = utils.NullInt32FromString(ctx.Query("age"))
	if err != nil {
		return age, gender, country, platform, offset, limit, err
	}
	if age.Valid && age.Int32 < 1 || age.Int32 > 100 {
		return age, gender, country, platform, offset, limit, NewInvalidQueryParameterError("age", "must be between 1 and 100")
	}
	gender = utils.NonEmptyNullStringFromString(ctx.Query("gender"))
	if gender.Valid {
		// 判斷 gender 在 db 中有資料
		count, err := app.databaseQueries.CheckGender(ctx, gender.String)
		if err != nil {
			return age, gender, country, platform, offset, limit, err
		}
		if count == 0 {
			return age, gender, country, platform, offset, limit, NewInvalidQueryParameterError("gender", "not in the database")
		}
	}
	country = utils.NonEmptyNullStringFromString(ctx.Query("country"))
	if country.Valid {
		// 判斷 country 在 db 中有資料
		count, err := app.databaseQueries.CheckCountry(ctx, country.String)
		if err != nil {
			return age, gender, country, platform, offset, limit, err
		}
		if count == 0 {
			return age, gender, country, platform, offset, limit, NewInvalidQueryParameterError("country", "not in the database")
		}
	}
	platform = utils.NonEmptyNullStringFromString(ctx.Query("platform"))
	if platform.Valid {
		// 判斷 platform 在 db 中有資料
		count, err := app.databaseQueries.CheckPlatform(ctx, platform.String)
		if err != nil {
			return age, gender, country, platform, offset, limit, err
		}
		if count == 0 {
			return age, gender, country, platform, offset, limit, NewInvalidQueryParameterError("platform", "not in the database")
		}
	}

	if offset_str := ctx.Query("offset"); offset_str != "" {
		offset, err = strconv.Atoi(offset_str)
		if err != nil {
			return age, gender, country, platform, offset, limit, err
		}
		if offset < 0 {
			return age, gender, country, platform, offset, limit, NewInvalidQueryParameterError("offset", "must be >= 0")
		}
	}

	if limit_str := ctx.Query("limit"); limit_str != "" {
		limit, err = strconv.Atoi(limit_str)
		if err != nil {
			return age, gender, country, platform, offset, limit, err
		}
		if limit < 0 {
			return age, gender, country, platform, offset, limit, NewInvalidQueryParameterError("limit", "must be >= 0")
		}
	}
	return age, gender, country, platform, offset, limit, nil
}

func (app *application) getAdvertisementHandler(ctx *gin.Context) {
	age, gender, country, platform, offset, limit, err := app.getAdvertisementQueryParameters(ctx)
	if err != nil {
		app.errorLogger.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ads, err := app.databaseQueries.GetAdvertisements(ctx, db.GetAdvertisementsParams{
		Age:      age,
		Gender:   gender,
		Country:  country,
		Platform: platform,
		Offset:   int32(offset),
		Limit:    int32(limit),
	})
	if err != nil {
		app.errorLogger.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": ads,
	})
}

func (app *application) createAdvertisementHandler(ctx *gin.Context) {
	body := CreateAdvertisementForm{}
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	advertisementId, err := app.databaseQueries.CreateAdvertisement(ctx, db.CreateAdvertisementParams{
		Title:   body.Title,
		StartAt: body.StartAt,
		EndAt:   body.EndAt,
	})
	if err != nil {
		app.errorLogger.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, condition := range body.Conditions {
		conditionId, err := app.databaseQueries.CreateCondition(ctx, db.CreateConditionParams{
			AgeStart: condition.AgeStart,
			AgeEnd:   condition.AgeEnd,
		})
		if err != nil {
			app.errorLogger.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		if condition.Gender.Valid {
			err = app.databaseQueries.CreateConditionGender(ctx, db.CreateConditionGenderParams{
				ConditionID: int32(conditionId),
				Gender:      condition.Gender.String,
			})
			if err != nil {
				app.errorLogger.Println(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
		if condition.Country.Valid {
			err = app.databaseQueries.CreateConditionCountry(ctx, db.CreateConditionCountryParams{
				ConditionID: int32(conditionId),
				Country:     condition.Country.String,
			})
			if err != nil {
				app.errorLogger.Println(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
		if condition.Platform.Valid {
			err = app.databaseQueries.CreateConditionPlatform(ctx, db.CreateConditionPlatformParams{
				ConditionID: int32(conditionId),
				Platform:    condition.Platform.String,
			})
			if err != nil {
				app.errorLogger.Println(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
		err = app.databaseQueries.CreateAdvertisementCondition(ctx, db.CreateAdvertisementConditionParams{
			AdvertisementID: int32(advertisementId),
			ConditionID:     int32(conditionId),
		})
		if err != nil {
			app.errorLogger.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
