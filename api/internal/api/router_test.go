package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kubmovie/api/internal/api"
	"kubmovie/api/internal/models"
	"kubmovie/api/internal/repository"
)

type fakeMovies struct {
	offset, limit int
	err           error
}

func (f *fakeMovies) List(_ context.Context, offset, limit int) ([]models.MovieListItem, error) {
	f.offset, f.limit = offset, limit
	return []models.MovieListItem{{MovieID: 1, Title: "Toy Story (1995)"}}, f.err
}

func (f *fakeMovies) GetByID(_ context.Context, id int) (*models.Movie, error) {
	if f.err != nil {
		return nil, f.err
	}
	if id != 1 {
		return nil, repository.ErrNotFound
	}
	return &models.Movie{MovieID: 1, Title: "Toy Story (1995)", Genres: "Adventure|Animation"}, nil
}

func do(t *testing.T, repo repository.MovieRepository, path string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	api.NewRouter(repo).ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
}

func TestHealth(t *testing.T) {
	w := do(t, &fakeMovies{}, "/health")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

func TestListMovies(t *testing.T) {
	f := &fakeMovies{}
	w := do(t, f, "/movies")
	require.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"pagenum":1,"movies":[{"movie_id":1,"title":"Toy Story (1995)"}]}`, w.Body.String())
	assert.Equal(t, 0, f.offset)
	assert.Equal(t, 50, f.limit)

	w = do(t, f, "/movies?pagenum=3")
	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 100, f.offset)
}

func TestListMoviesBadPagenum(t *testing.T) {
	for _, q := range []string{"0", "-1", "abc"} {
		w := do(t, &fakeMovies{}, "/movies?pagenum="+q)
		assert.Equal(t, http.StatusBadRequest, w.Code, q)
	}
}

func TestListMoviesError(t *testing.T) {
	w := do(t, &fakeMovies{err: errors.New("boom")}, "/movies")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetMovie(t *testing.T) {
	w := do(t, &fakeMovies{}, "/movies/1")
	require.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"movie_id":1,"title":"Toy Story (1995)","genres":"Adventure|Animation"}`, w.Body.String())
}

func TestGetMovieNotFound(t *testing.T) {
	w := do(t, &fakeMovies{}, "/movies/999")
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.JSONEq(t, `{"error":"movie not found"}`, w.Body.String())
}

func TestGetMovieBadID(t *testing.T) {
	w := do(t, &fakeMovies{}, "/movies/abc")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
