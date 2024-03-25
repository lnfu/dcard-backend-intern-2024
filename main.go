package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.GET("/", hello)

	apiV1 := router.Group("api/v1/")
	apiV1.POST("ad", createAdvertisement)
	apiV1.GET("ad", listAdvertisement)

	router.Run(":8080")
}

func hello(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Hello, World!")
}

type Advertisement struct {
	Title      string      `json:"title" binding:"required"`
	AtartAt    time.Time   `json:"startAt" binding:"required"`
	EndAt      time.Time   `json:"endAt" binding:"required"`
	Conditions []Condition `json:"conditions"`
}

type Condition struct {
	AgeStart int      `json:"ageStart"`
	AgeEnd   int      `json:"ageEnd"`
	Country  []string `json:"country"`
	Platform []string `json:"platform"`
	Gender   []string `json:"gender"`
}

func listAdvertisement(ctx *gin.Context) {
	offset, _ := strconv.Atoi(ctx.Query("offset"))
	limit, _ := strconv.Atoi(ctx.Query("limit"))

	if limit == 0 {
		limit = 5
	}

	age := strings.Split(ctx.Query("age"), ",")
	gender := strings.Split(ctx.Query("gender"), ",")
	country := strings.Split(ctx.Query("country"), ",")
	platform := strings.Split(ctx.Query("platform"), ",")
	_, _, _, _, _ = offset, age, gender, country, platform

	// TODO SQL SELECT

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func createAdvertisement(ctx *gin.Context) {
	advertisement := Advertisement{}
	fmt.Println(ctx.Request)

	if err := ctx.BindJSON(&advertisement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO SQL INSERT INTO

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
