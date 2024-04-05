package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/lnfu/dcard-intern/db/sqlc"

	"github.com/lnfu/dcard-intern/utils"
)

func (handler *Handler) getAdvertisementQueryParameters(ctx *gin.Context) (sql.NullInt32, sql.NullString, sql.NullString, sql.NullString, int, int, error) {
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
		// validate gender
		if !handler.genderSet.Contains(gender.String) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid gender value"})
			return age, gender, country, platform, offset, limit, errors.New("")
		}
	}
	country = utils.NonEmptyNullStringFromString(ctx.Query("country"))
	if country.Valid {
		// validate country
		if !handler.countrySet.Contains(country.String) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country value"})
			return age, gender, country, platform, offset, limit, errors.New("")
		}
	}
	platform = utils.NonEmptyNullStringFromString(ctx.Query("platform"))
	if platform.Valid {
		// validate platform
		if !handler.platformSet.Contains(platform.String) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid platform value"})
			return age, gender, country, platform, offset, limit, errors.New("")
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
func (handler *Handler) GetAdvertisementHandler(ctx *gin.Context) {
	age, gender, country, platform, offset, limit, err := handler.getAdvertisementQueryParameters(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	params := db.GetActiveAdvertisementsParams{
		Age:      age,
		Gender:   gender,
		Country:  country,
		Platform: platform,
		Offset:   int32(offset),
		Limit:    int32(limit),
	}

	// search cache
	ads, err := handler.cac.GetAdvertisementsFromCache(ctx, params)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"items": ads,
		})
		return
	}

	// search database
	ads, err = handler.databaseQueries.GetActiveAdvertisements(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = handler.cac.SetAdvertisementsToCache(ctx, params, ads)
	if err != nil {
		log.Fatalf("Redis: 無法新增 Cache (%v)\n", err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"items": ads,
	})
}
