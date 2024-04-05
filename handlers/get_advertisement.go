package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/lnfu/dcard-intern/db/sqlc"
)

type QueryParameters struct {
	Age      int32  `form:"age" example:"24"`
	Gender   string `form:"gender" example:"M"`
	Country  string `form:"country" example:"TW"`
	Platform string `form:"platform" example:"android"`
	Offset   int32  `form:"offset" example:"0"`
	Limit    int32  `form:"limit" example:"5"`
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

	var params db.GetActiveAdvertisementsParams

	// age
	if queryParameters.Age < 0 || queryParameters.Age > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age value (must be 1 ~ 100)"})
		return
	} else if queryParameters.Age == 0 {
		params.Age = sql.NullInt32{Int32: 0, Valid: false}
	} else {
		params.Age = sql.NullInt32{Int32: queryParameters.Age, Valid: true}
	}

	// gender
	if queryParameters.Gender == "" {
		params.Gender = sql.NullString{String: "", Valid: false}
	} else if !handler.genderSet.Contains(queryParameters.Gender) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender value"})
		return
	} else {
		params.Gender = sql.NullString{String: queryParameters.Gender, Valid: true}
	}

	// country
	if queryParameters.Country == "" {
		params.Country = sql.NullString{String: "", Valid: false}
	} else if !handler.countrySet.Contains(queryParameters.Country) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country value"})
		return
	} else {
		params.Country = sql.NullString{String: queryParameters.Country, Valid: true}
	}

	// platform
	if queryParameters.Platform == "" {
		params.Platform = sql.NullString{String: "", Valid: false}
	} else if !handler.platformSet.Contains(queryParameters.Platform) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid platform value"})
		return
	} else {
		params.Platform = sql.NullString{String: queryParameters.Platform, Valid: true}
	}

	// offset
	params.Offset = queryParameters.Offset

	// limit
	if queryParameters.Limit == 0 {
		params.Limit = 5
	} else if queryParameters.Limit < 1 || queryParameters.Limit > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value (must be 1 ~ 100)"})
		return
	} else {
		params.Limit = queryParameters.Limit
	}

	// find in cache
	ads, err := handler.cac.GetAdvertisementsFromCache(ctx, params)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"items": ads,
		})
		return
	}

	// find in database
	ads, err = handler.databaseQueries.GetActiveAdvertisements(ctx, params)
	if err != nil {
		log.Println("Database error: :))))", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	// add cache
	err = handler.cac.SetAdvertisementsToCache(ctx, params, ads)
	if err != nil {
		log.Println("Cache error:", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cache error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": ads,
	})
}
