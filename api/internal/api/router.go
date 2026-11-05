package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"kubmovie/api/internal/repository"
)

const PageSize = 50

type handler struct {
	movies repository.MovieRepository
}

func NewRouter(movies repository.MovieRepository) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	h := &handler{movies: movies}
	r.GET("/health", h.health)
	r.GET("/movies", h.listMovies)
	r.GET("/movies/:id", h.getMovie)
	return r
}

func (h *handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *handler) listMovies(c *gin.Context) {
	pagenum := 1
	if s := c.Query("pagenum"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pagenum must be an integer >= 1"})
			return
		}
		pagenum = n
	}

	items, err := h.movies.List(c.Request.Context(), (pagenum-1)*PageSize, PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pagenum": pagenum, "movies": items})
}

func (h *handler) getMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	m, err := h.movies.GetByID(c.Request.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, m)
}
