package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"kubmovie/api/internal/config"
	"kubmovie/api/internal/database"
	"kubmovie/api/internal/repository"
)

// These tests run against the Postgres database configured in ../../.env.

var testDB *gorm.DB

func TestMain(m *testing.M) {
	cfg, err := config.LoadFile("../../.env")
	if err != nil {
		panic(err)
	}
	testDB, err = database.Connect(cfg.DSN())
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestConnectivity(t *testing.T) {
	sqlDB, err := testDB.DB()
	require.NoError(t, err)
	assert.NoError(t, sqlDB.PingContext(context.Background()))

	var one int
	require.NoError(t, testDB.Raw("SELECT 1").Scan(&one).Error)
	assert.Equal(t, 1, one)
}

func TestMovieList(t *testing.T) {
	repo := repository.NewMovieRepository(testDB)

	page1, err := repo.List(context.Background(), 0, 50)
	require.NoError(t, err)
	assert.Len(t, page1, 50)

	page2, err := repo.List(context.Background(), 50, 50)
	require.NoError(t, err)
	assert.Len(t, page2, 50)

	assert.Less(t, page1[len(page1)-1].MovieID, page2[0].MovieID, "results are ordered by movie_id with no overlap")
	assert.Equal(t, 1, page1[0].MovieID)
	assert.NotEmpty(t, page1[0].Title)
}

func TestMovieGetByID(t *testing.T) {
	repo := repository.NewMovieRepository(testDB)

	m, err := repo.GetByID(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, m.MovieID)
	assert.Contains(t, m.Title, "Toy Story")
	assert.Contains(t, m.Genres, "Animation")

	_, err = repo.GetByID(context.Background(), -1)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestLinkGetByMovieID(t *testing.T) {
	repo := repository.NewLinkRepository(testDB)

	l, err := repo.GetByMovieID(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, 1, l.MovieID)
	assert.NotEmpty(t, l.ImdbID)

	_, err = repo.GetByMovieID(context.Background(), -1)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestRatingListByMovieID(t *testing.T) {
	repo := repository.NewRatingRepository(testDB)

	ratings, err := repo.ListByMovieID(context.Background(), 1, 0, 10)
	require.NoError(t, err)
	require.NotEmpty(t, ratings)
	assert.LessOrEqual(t, len(ratings), 10)
	for _, r := range ratings {
		assert.Equal(t, 1, r.MovieID)
		assert.GreaterOrEqual(t, r.Rating, 0.5)
		assert.LessOrEqual(t, r.Rating, 5.0)
	}
}

func TestTagListByMovieID(t *testing.T) {
	repo := repository.NewTagRepository(testDB)

	var movieID int
	require.NoError(t, testDB.Raw("SELECT movie_id FROM tags LIMIT 1").Scan(&movieID).Error)
	require.NotZero(t, movieID, "tags table should contain data")

	tags, err := repo.ListByMovieID(context.Background(), movieID)
	require.NoError(t, err)
	require.NotEmpty(t, tags)
	for _, tg := range tags {
		assert.Equal(t, movieID, tg.MovieID)
		assert.NotEmpty(t, tg.Tag)
	}
}
