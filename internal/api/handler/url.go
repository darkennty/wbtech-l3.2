package handler

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"WBTech_L3.2/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type UrlRequest struct {
	LongUrl         string `json:"url"`
	DesiredShortUrl string `json:"short_url,omitempty"`
}

func (h *Handler) handleCreate(c *gin.Context) {
	var request UrlRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		ReturnErrorResponse(c, http.StatusBadRequest, "invalid json")
		return
	}

	longUrl, err := url.Parse(request.LongUrl)
	if err != nil {
		ReturnErrorResponse(c, http.StatusBadRequest, "invalid url")
		return
	}

	if longUrl.Scheme == "" {
		longUrl.Scheme = "https"
	}

	if _, err = http.Get(longUrl.String()); err != nil {
		ReturnErrorResponse(c, http.StatusBadRequest, "non-existing url")
		return
	}

	if request.DesiredShortUrl != "" {
		if len(request.DesiredShortUrl) < 5 {
			ReturnErrorResponse(c, http.StatusBadRequest, "short_url is too short (need at least 5 characters)")
			return
		}

		if len(request.DesiredShortUrl) > 32 {
			ReturnErrorResponse(c, http.StatusBadRequest, "short_url is too long (should be at most 32 characters)")
			return
		}

		if !parseCustomShortUrl(request.DesiredShortUrl) {
			ReturnErrorResponse(c, http.StatusBadRequest, "short_url contains unallowed characters")
			return
		}
	}

	shortUrl, err := h.services.CreateShortUrl(c, longUrl.String(), request.DesiredShortUrl)
	if err != nil {
		ReturnErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnResultResponse(c, http.StatusOK, ginext.H{"short_url": shortUrl})
}

func (h *Handler) handleRedirect(c *gin.Context) {
	shortUrl := c.Param("short_url")
	if shortUrl == "" {
		ReturnErrorResponse(c, http.StatusBadRequest, "no url given")
		return
	}

	longUrl, err := h.services.GetLongUrl(c, shortUrl)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			ReturnErrorResponse(c, http.StatusNotFound, "404 not found")
			return
		}
		ReturnErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, longUrl)

	userAgent := c.GetHeader("User-Agent")
	err = h.services.SaveStats(c, shortUrl, userAgent)
	if err != nil {
		ReturnErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) handleGetStats(c *gin.Context) {
	shortUrl := c.Param("short_url")
	if shortUrl == "" {
		ReturnErrorResponse(c, http.StatusBadRequest, "no url given")
		return
	}

	var stat any
	var err error
	aggregateBy := c.Query("aggregate_by")
	aggregateBy = strings.ToLower(aggregateBy)

	if aggregateBy == "day" || aggregateBy == "month" || aggregateBy == "user_agent" || aggregateBy == "useragent" {
		stat, err = h.services.GetAggregatedStats(c, shortUrl, aggregateBy)
		if err != nil {
			ReturnErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		ReturnResultResponse(c, http.StatusOK, ginext.H{"aggregated_by": aggregateBy, "stats": stat})
	} else {
		stat, err = h.services.GetStats(c, shortUrl)
		if err != nil {
			ReturnErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		ReturnResultResponse(c, http.StatusOK, ginext.H{"stats": stat})
	}
}

func parseCustomShortUrl(shortUrl string) bool {
	pattern := regexp.MustCompile("^[A-Za-z0-9_-]+$")
	return pattern.MatchString(shortUrl)
}
