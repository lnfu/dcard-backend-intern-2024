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
	Title      string                   `json:"title" example:"AD 55" extensions:"x-order=0"`
	StartAt    time.Time                `json:"startAt" example:"2023-12-10T03:00:00.000Z" extensions:"x-order=1"`
	EndAt      time.Time                `json:"endAt" example:"2023-12-31T16:00:00.000Z" extensions:"x-order=2"`
	Conditions []AdvertisementCondition `json:"conditions" extensions:"x-order=3"`
}

type AdvertisementCondition struct {
	// TODO gender, country, platform 可以多選
	AgeStart *int    `json:"ageStart,omitempty" example:"20" swaggertype:"integer" extensions:"x-order=0"`
	AgeEnd   *int    `json:"ageEnd,omitempty" example:"30" swaggertype:"integer" extensions:"x-order=1"`
	Gender   *string `json:"gender,omitempty" example:"M" swaggertype:"string" extensions:"x-order=2"`
	Country  *string `json:"country,omitempty" example:"TW" swaggertype:"string" extensions:"x-order=3"`
	Platform *string `json:"platform,omitempty" example:"ios" swaggertype:"string" extensions:"x-order=4"`
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
	// TODO 關於 database 的錯誤理論上應該是回 internal server error
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

// @Summary		列出符合可⽤和匹配⽬標條件的廣告
// @BasePath	/api/v1
// @Version		1.0
// @Param		age query int false "年齡條件" minimum(1) maximum(100)
// @Param		gender query string false "性別條件 (M/F)" Enums(M, F)
// @Param		country query string false "國家條件 (參考 ISO_3166-1 alpha-2)"
// @Param		platform query string false "平台條件" Enums(android, ios, web)
// @Param		offset query int false " "
// @Param		limit query int false " "
// @Produce		json
// @Tags		advertisement
// @Router		/ad [get]
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

// @Summary		產⽣廣告資源
// @BasePath	/api/v1
// @Version		1.0
// @Param		request body main.CreateAdvertisementForm true "廣告內容"
// @Produce		json
// @Tags		advertisement
// @Router		/ad [post]
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
			AgeStart: utils.NullInt32FromInt32Pointer(condition.AgeStart),
			AgeEnd:   utils.NullInt32FromInt32Pointer(condition.AgeEnd),
		})
		if err != nil {
			app.errorLogger.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		if condition.Gender != nil {
			err = app.databaseQueries.CreateConditionGender(ctx, db.CreateConditionGenderParams{
				ConditionID: int32(conditionId),
				Gender:      *(condition.Gender),
			})
			if err != nil {
				app.errorLogger.Println(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
		if condition.Country != nil {
			err = app.databaseQueries.CreateConditionCountry(ctx, db.CreateConditionCountryParams{
				ConditionID: int32(conditionId),
				Country:     *(condition.Country),
			})
			if err != nil {
				app.errorLogger.Println(err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
		if condition.Platform != nil {
			err = app.databaseQueries.CreateConditionPlatform(ctx, db.CreateConditionPlatformParams{
				ConditionID: int32(conditionId),
				Platform:    *(condition.Platform),
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
