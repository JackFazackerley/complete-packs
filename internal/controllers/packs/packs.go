package packs

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/JackFazackerley/complete-packs/pkg/cache"

	"github.com/JackFazackerley/complete-packs/internal/database"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	db         database.Database
	sizesCache cache.Cache
}

type Request struct {
	Size int `json:"size"`
}

func New(db database.Database, sizesCache cache.Cache) *Controller {
	return &Controller{db: db, sizesCache: sizesCache}
}

func (r *Controller) Read(c *gin.Context) {
	packs, err := r.db.ReadPacks()
	if err != nil {
		log.WithError(err).Error("reading packs")
		c.JSON(http.StatusInternalServerError, c.Error(errors.Cause(err)))
		return
	}

	c.JSON(http.StatusOK, packs)
}

func (r *Controller) Write(c *gin.Context) {
	var input Request
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithError(err).Error("parsing request body")
		c.Status(http.StatusBadRequest)
		return
	}

	if err := r.db.WritePack(input.Size); err != nil {
		log.WithError(err).Error("write new size")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := r.sizesCache.Load(); err != nil {
		log.WithError(err).Error("loading cache")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (r *Controller) Delete(c *gin.Context) {
	var input Request
	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithError(err).Error("parsing request body")
		c.Status(http.StatusBadRequest)
		return
	}

	if err := r.db.DeletePack(input.Size); err != nil {
		log.WithError(err).Error("deleting size")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := r.sizesCache.Load(); err != nil {
		log.WithError(err).Error("loading cache")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
