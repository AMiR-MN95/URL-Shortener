package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type ShortURL struct {
	ID       string `json:"id"`
	Original string `json:"original"`
}

var urlMap = make(map[string]ShortURL)

// Generate Short URLs

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortURL() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Initialize Gin

func main() {
	router := gin.Default()

	router.POST("/shorten", shortenURL)
	router.GET("/:shortURL", redirectURL)

	router.Run(":8080")
}

// Shorten URL Handler
func shortenURL(c *gin.Context) {
	var requestBody struct {
		Original string `json:"original"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	shortURL := generateShortURL()
	urlMap[shortURL] = ShortURL{ID: shortURL, Original: requestBody.Original}

	c.JSON(http.StatusOK, gin.H{"short.url": shortURL})
}

// Redirect URL Handler
func redirectURL(c *gin.Context) {
	shortURL := c.Param("shortURL")
	url, exists := urlMap[shortURL]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
	}

	c.Redirect(http.StatusFound, url.Original)
}
