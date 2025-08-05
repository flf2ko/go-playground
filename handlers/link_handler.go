package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/flf2ko/playground/go-api-sample/database"
	"github.com/flf2ko/playground/go-api-sample/models"
	"github.com/flf2ko/playground/go-api-sample/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type LinkHandler struct {
	db     *database.DB
	client *resty.Client
}

func NewLinkHandler(db *database.DB) *LinkHandler {
	client := resty.New()
	client.SetTimeout(30 * time.Second)
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)

	return &LinkHandler{
		db:     db,
		client: client,
	}
}

func (h *LinkHandler) FetchJSON(c *gin.Context) {
	link := c.Query("link")
	if link == "" {
		c.JSON(http.StatusBadRequest, models.FetchResponse{
			Success: false,
			Message: "Missing required parameter: link",
			Error:   "link parameter is required",
		})
		return
	}

	if err := utils.IsValidURL(link); err != nil {
		c.JSON(http.StatusBadRequest, models.FetchResponse{
			Success: false,
			Message: "Invalid URL format",
			Error:   err.Error(),
		})
		return
	}

	log.Printf("Fetching JSON from URL: %s", link)

	resp, err := h.client.R().
		SetContext(c).
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "go-api-sample/1.0").
		Get(link)

	if err != nil {
		log.Printf("Failed to fetch URL %s: %v", link, err)
		c.JSON(http.StatusBadRequest, models.FetchResponse{
			Success: false,
			Message: "Failed to fetch URL",
			Error:   fmt.Sprintf("HTTP request failed: %v", err),
		})
		return
	}

	if resp.StatusCode() >= 400 {
		log.Printf("HTTP error for URL %s: %d %s", link, resp.StatusCode(), resp.Status())
		c.JSON(http.StatusBadRequest, models.FetchResponse{
			Success: false,
			Message: "HTTP request failed",
			Error:   fmt.Sprintf("HTTP %d: %s", resp.StatusCode(), resp.Status()),
		})
		return
	}

	contentType := resp.Header().Get("Content-Type")
	if !utils.IsJSONContentType(contentType) {
		log.Printf("Invalid content type for URL %s: %s", link, contentType)
		c.JSON(http.StatusBadRequest, models.FetchResponse{
			Success: false,
			Message: "Response is not JSON",
			Error:   fmt.Sprintf("Expected JSON content, got: %s", contentType),
		})
		return
	}

	body := resp.String()
	if err := utils.IsValidJSON(body); err != nil {
		log.Printf("Invalid JSON from URL %s: %v", link, err)
		c.JSON(http.StatusBadRequest, models.FetchResponse{
			Success: false,
			Message: "Invalid JSON response",
			Error:   err.Error(),
		})
		return
	}

	record, err := h.db.SaveJSONRecord(c.Request.Context(), link, body)
	if err != nil {
		log.Printf("Failed to save JSON record for URL %s: %v", link, err)
		c.JSON(http.StatusInternalServerError, models.FetchResponse{
			Success: false,
			Message: "Failed to save to database",
			Error:   "Internal server error",
		})
		return
	}

	log.Printf("Successfully saved JSON record ID: %d for URL: %s", record.ID, link)

	c.JSON(http.StatusOK, models.FetchResponse{
		Success: true,
		Message: "JSON fetched and saved successfully",
		Data:    record,
	})
}

func (h *LinkHandler) GetRecords(c *gin.Context) {
	records, err := h.db.GetJSONRecords(c.Request.Context(), 10)
	if err != nil {
		log.Printf("Failed to get records: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve records",
			"error":   "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Records retrieved successfully",
		"data":    records,
		"count":   len(records),
	})
}
