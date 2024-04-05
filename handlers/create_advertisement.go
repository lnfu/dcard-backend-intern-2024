package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/lnfu/dcard-intern/db/sqlc"

	"github.com/lnfu/dcard-intern/utils"
)

type Advertisement struct {
	Title      string                   `json:"title" binding:"required" example:"AD 55" extensions:"x-order=0"`
	StartAt    time.Time                `json:"startAt" binding:"required" example:"2023-12-10T03:00:00.000Z" extensions:"x-order=1"`
	EndAt      time.Time                `json:"endAt" binding:"required" example:"2023-12-31T16:00:00.000Z" extensions:"x-order=2"`
	Conditions []AdvertisementCondition `json:"conditions" extensions:"x-order=3"`
}

type AdvertisementCondition struct {
	AgeStart *int     `json:"ageStart,omitempty" example:"20" swaggertype:"integer" extensions:"x-order=0"`
	AgeEnd   *int     `json:"ageEnd,omitempty" example:"30" swaggertype:"integer" extensions:"x-order=1"`
	Gender   []string `json:"gender,omitempty" example:"M" swaggertype:"array,string" extensions:"x-order=2"`
	Country  []string `json:"country,omitempty" example:"TW,JP" swaggertype:"array,string" extensions:"x-order=3"`
	Platform []string `json:"platform,omitempty" example:"android,ios" swaggertype:"array,string" extensions:"x-order=4"`
}

func NewInvalidQueryParameterError(parameterName, reason string) InvalidQueryParameterError {
	return InvalidQueryParameterError{
		ParameterName: parameterName,
		Reason:        reason,
	}
}

// @Summary		產⽣廣告資源
// @BasePath	/api/v1
// @Version		1.0
// @Param		request body handlers.CreateAdvertisementForm true "廣告內容"
// @Produce		json
// @Tags		advertisement
// @Router		/ad [post]
func (handler *Handler) CreateAdvertisementHandler(ctx *gin.Context) {
	// get params from request body
	body := Advertisement{}
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// add ad to database
	advertisementId, err := handler.databaseQueries.CreateAdvertisement(ctx, db.CreateAdvertisementParams{
		Title:   body.Title,
		StartAt: body.StartAt,
		EndAt:   body.EndAt,
	})
	if err != nil {
		log.Println("Database error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	for _, condition := range body.Conditions {
		// 判斷年齡是 1 ~ 100
		ageStart := utils.NullInt32FromInt32Pointer(condition.AgeStart)
		if ageStart.Valid && ageStart.Int32 < 1 || ageStart.Int32 > 100 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": NewInvalidQueryParameterError("ageStart", "must be between 1 and 100").Error()})
			return
		}
		ageEnd := utils.NullInt32FromInt32Pointer(condition.AgeEnd)
		if ageEnd.Valid && ageEnd.Int32 < 1 || ageEnd.Int32 > 100 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": NewInvalidQueryParameterError("ageEnd", "must be between 1 and 100").Error()})
			return
		}
		// 判斷 ageStart <= ageEnd
		if ageStart.Valid && ageEnd.Valid && ageStart.Int32 > ageEnd.Int32 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": NewInvalidQueryParameterError("ageEnd", "must be greater than ageStart").Error()})
			return
		}
		conditionId, err := handler.databaseQueries.CreateCondition(ctx, db.CreateConditionParams{
			AgeStart: ageStart,
			AgeEnd:   utils.NullInt32FromInt32Pointer(condition.AgeEnd),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, gender := range condition.Gender {
			// validate gender
			if !handler.genderSet.Contains(gender) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender value"})
				return
			}

			err = handler.databaseQueries.CreateConditionGender(ctx, db.CreateConditionGenderParams{
				ConditionID: int32(conditionId),
				Gender:      gender,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}

		for _, country := range condition.Country {
			// validate country
			if !handler.countrySet.Contains(country) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country value"})
				return
			}

			err = handler.databaseQueries.CreateConditionCountry(ctx, db.CreateConditionCountryParams{
				ConditionID: int32(conditionId),
				Country:     country,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}

		for _, platform := range condition.Platform {
			// 判斷 platform 在 cache/db 中有資料
			if !handler.platformSet.Contains(platform) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid platform value"})
				return
			}

			err = handler.databaseQueries.CreateConditionPlatform(ctx, db.CreateConditionPlatformParams{
				ConditionID: int32(conditionId),
				Platform:    platform,
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
		}
		err = handler.databaseQueries.CreateAdvertisementCondition(ctx, db.CreateAdvertisementConditionParams{
			AdvertisementID: int32(advertisementId),
			ConditionID:     int32(conditionId),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
