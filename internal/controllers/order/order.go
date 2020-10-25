package order

import (
	"errors"
	"net/http"

	"github.com/JackFazackerley/complete-packs/pkg/cache"
	"github.com/JackFazackerley/complete-packs/pkg/order/best"
	"github.com/JackFazackerley/complete-packs/pkg/order/fast"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	parsingRequestBody = errors.New("parsing request body")
	targetTooLarge     = errors.New("target too large")
)

type Controller struct {
	sizeCache cache.Cache
}

type Request struct {
	Target float64 `json:"target"`
}

func New(sizeCache cache.Cache) *Controller {
	return &Controller{
		sizeCache: sizeCache,
	}
}

func (o *Controller) Best(c *gin.Context) {
	var input Request

	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithError(err).Error(parsingRequestBody)
		c.JSON(http.StatusBadRequest, c.Error(parsingRequestBody))
		return
	}

	// this is purely here to stop the server crashing and running out of memory
	if input.Target > 9999999 {
		c.JSON(http.StatusBadRequest, c.Error(targetTooLarge))
		return
	}

	sizes := o.sizeCache.Asc()

	result := best.Calculate(int(input.Target), sizes)

	c.JSON(http.StatusOK, result)
}

func (o *Controller) Fast(c *gin.Context) {
	var input Request
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithError(err).Error(parsingRequestBody)
		c.JSON(http.StatusBadRequest, c.Error(parsingRequestBody))
		return
	}

	sizes := o.sizeCache.Desc()

	result := fast.Calculate(input.Target, sizes)

	c.JSON(http.StatusOK, result)
}
