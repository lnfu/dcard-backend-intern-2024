package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sqlc "github.com/lnfu/dcard-intern/app/models/sqlc"

	"github.com/lnfu/dcard-intern/app/utils"
)

type Advertisement struct {
	Title      string                   `json:"title" binding:"required" example:"AD 55" extensions:"x-order=0"`
	StartAt    time.Time                `json:"startAt" binding:"required" example:"2023-12-10T03:00:00.000Z" extensions:"x-order=1"`
	EndAt      time.Time                `json:"endAt" binding:"required" example:"2023-12-31T16:00:00.000Z" extensions:"x-order=2"`
	Conditions []AdvertisementCondition `json:"conditions" extensions:"x-order=3"`
}

type AdvertisementCondition struct {
	AgeStart *int32   `json:"ageStart,omitempty" example:"20" swaggertype:"integer" extensions:"x-order=0"`
	AgeEnd   *int32   `json:"ageEnd,omitempty" example:"30" swaggertype:"integer" extensions:"x-order=1"`
	Gender   []string `json:"gender,omitempty" example:"M" swaggertype:"array,string" extensions:"x-order=2"`
	Country  []string `json:"country,omitempty" example:"TW,JP" swaggertype:"array,string" extensions:"x-order=3"`
	Platform []string `json:"platform,omitempty" example:"android,ios" swaggertype:"array,string" extensions:"x-order=4"`
}

// @Summary		產⽣廣告資源
// @BasePath	/api/v1
// @Version		1.0
// @Param		request body handlers.Advertisement true "廣告內容"
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

	// TODO validate startAt <= endAt

	// add ad to database
	advertisementId, err := handler.databaseQueries.CreateAdvertisement(ctx, sqlc.CreateAdvertisementParams{
		Title:   body.Title,
		StartAt: body.StartAt,
		EndAt:   body.EndAt,
	})
	if err != nil {
		log.Println("Database error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// process conditions
	for _, condition := range body.Conditions {
		if err := handler.validateCondition(condition); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// add condition
		conditionId, err := handler.databaseQueries.CreateCondition(ctx, sqlc.CreateConditionParams{
			AgeStart: utils.NullInt32FromInt32Pointer(condition.AgeStart),
			AgeEnd:   utils.NullInt32FromInt32Pointer(condition.AgeEnd),
		})

		if err != nil {
			log.Println("Database Error:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		for _, gender := range condition.Gender {

			// add gender-condition relation
			err = handler.databaseQueries.CreateConditionGender(ctx, sqlc.CreateConditionGenderParams{
				ConditionID: int32(conditionId),
				Gender:      gender,
			})
			if err != nil {
				log.Println("Database Error:", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
				return
			}
		}

		for _, country := range condition.Country {

			// add country-condition relation
			err = handler.databaseQueries.CreateConditionCountry(ctx, sqlc.CreateConditionCountryParams{
				ConditionID: int32(conditionId),
				Country:     country,
			})
			if err != nil {
				log.Println("Database Error:", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
				return
			}
		}

		for _, platform := range condition.Platform {

			// add platform-condition relation
			err = handler.databaseQueries.CreateConditionPlatform(ctx, sqlc.CreateConditionPlatformParams{
				ConditionID: int32(conditionId),
				Platform:    platform,
			})
			if err != nil {
				log.Println("Database Error:", err.Error())
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
				return
			}
		}

		// add condition-advertisement relation
		err = handler.databaseQueries.CreateAdvertisementCondition(ctx, sqlc.CreateAdvertisementConditionParams{
			AdvertisementID: int32(advertisementId),
			ConditionID:     int32(conditionId),
		})
		if err != nil {
			log.Println("Database Error:", err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (handler *Handler) validateCondition(condition AdvertisementCondition) error {
	// ageStart
	if condition.AgeStart != nil && (*condition.AgeStart < 1 || *condition.AgeStart > 100) {
		return errors.New("invalid ageStart value (must be 1 ~ 100)")
	}

	// ageEnd
	if condition.AgeEnd != nil && (*condition.AgeEnd < 1 || *condition.AgeEnd > 100) {
		return errors.New("invalid ageEnd value (must be 1 ~ 100)")
	}

	// ageStart <= ageEnd
	if condition.AgeStart != nil && condition.AgeEnd != nil && *condition.AgeStart > *condition.AgeEnd {
		return errors.New("invalid ageEnd value (must be >= ageStart)")
	}

	// gender
	for _, gender := range condition.Gender {
		if !handler.genderSet.Contains(gender) {
			return errors.New("invalid gender value")
		}
	}

	// country
	for _, country := range condition.Country {
		if !handler.countrySet.Contains(country) {
			return errors.New("invalid country value")
		}
	}

	// platform
	for _, platform := range condition.Platform {
		if !handler.platformSet.Contains(platform) {
			return errors.New("invalid platform value")
		}
	}

	return nil
}
