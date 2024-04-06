package handlers

import (
	"errors"
	"log"
	"net/http"

	db "github.com/lnfu/dcard-intern/app/models/sqlc"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
	"github.com/lnfu/dcard-intern/app/utils"
)

type QueryParameters struct {
	Age      *int32  `form:"age" example:"24"`
	Gender   *string `form:"gender" example:"M"`
	Country  *string `form:"country" example:"TW"`
	Platform *string `form:"platform" example:"android"`
	Offset   *int32  `form:"offset" example:"0"`
	Limit    *int32  `form:"limit" example:"5"`
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
func (handler *Handler) GetAdvertisementHandler(ctx *gin.Context) {
	var queryParameters QueryParameters
	if err := ctx.ShouldBindQuery(&queryParameters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := handler.validateQueryParameters(queryParameters); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	params := handler.buildDBParams(queryParameters)

	ads, err := handler.retrieveAdvertisements(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": ads,
	})
}

// 判斷 query parameters 是否 valid
func (handler *Handler) validateQueryParameters(queryParameters QueryParameters) error {

	// age
	if queryParameters.Age != nil && (*queryParameters.Age < 1 || *queryParameters.Age > 100) {
		return errors.New("invalid age value (must be 1 ~ 100)")
	}

	// gender
	if queryParameters.Gender != nil && !handler.genderSet.Contains(*queryParameters.Gender) {
		return errors.New("invalid gender value")
	}

	// country
	if queryParameters.Country != nil && !handler.countrySet.Contains(*queryParameters.Country) {
		return errors.New("invalid country value")
	}

	// platform
	if queryParameters.Platform != nil && !handler.platformSet.Contains(*queryParameters.Platform) {
		return errors.New("invalid platform value")
	}

	// offset
	if queryParameters.Offset != nil && (*queryParameters.Offset < 0) {
		return errors.New("invalid offset value (must be >= 0)")
	}

	// limit
	if queryParameters.Limit != nil && (*queryParameters.Limit < 1 || *queryParameters.Limit > 100) {
		return errors.New("invalid limit value (must be 1 ~ 100)")
	}

	return nil
}

// query parameters -> db parameters (如果是空的設定 null)
func (handler *Handler) buildDBParams(queryParameters QueryParameters) db.GetActiveAdvertisementsParams {
	var params db.GetActiveAdvertisementsParams

	// age
	params.Age = utils.NullInt32FromInt32Pointer(queryParameters.Age)

	// gender/country/platform
	params.Gender = utils.NullStringFromStringPointer(queryParameters.Gender)
	params.Country = utils.NullStringFromStringPointer(queryParameters.Country)
	params.Platform = utils.NullStringFromStringPointer(queryParameters.Platform)

	// offset
	if queryParameters.Offset == nil {
		params.Offset = 0
	} else {
		params.Offset = *queryParameters.Offset
	}

	// limit
	if queryParameters.Limit == nil {
		params.Limit = 5
	} else {
		params.Limit = *queryParameters.Limit
	}

	return params
}

// 從 cache/database 獲取符合條件的 advertisement
func (handler *Handler) retrieveAdvertisements(params db.GetActiveAdvertisementsParams) ([]db.Advertisement, error) {
	var ads []db.Advertisement

	// find in cache
	ads, err := handler.cac.GetAdvertisementsFromCache(ctx, params)
	if err == redis.Nil {
		// 沒找到, 去 database 找
		ads, err = handler.databaseQueries.GetActiveAdvertisements(ctx, params)
		if err != nil {
			log.Println("Database Error: ", err.Error())
			return nil, errors.New("database error")
		}

		// add cache
		err = handler.cac.SetAdvertisementsToCache(ctx, params, ads)
		if err != nil {
			log.Println("Cache Error: ", err.Error())
			return nil, errors.New("cache error")
		}

	} else if err != nil {
		// redis error
		log.Println("Cache Error: ", err.Error())
		return nil, errors.New("cache error")

	}

	return ads, nil
}
