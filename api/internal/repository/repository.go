package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"kubmovie/api/internal/models"
)

// ErrNotFound is returned when a requested row does not exist.
var ErrNotFound = errors.New("not found")

// Read-only repositories: no Create/Update/Delete.

type MovieRepository interface {
	List(ctx context.Context, offset, limit int) ([]models.MovieListItem, error)
	GetByID(ctx context.Context, id int) (*models.Movie, error)
}

type LinkRepository interface {
	GetByMovieID(ctx context.Context, movieID int) (*models.Link, error)
}

type RatingRepository interface {
	ListByMovieID(ctx context.Context, movieID, offset, limit int) ([]models.Rating, error)
}

type TagRepository interface {
	ListByMovieID(ctx context.Context, movieID int) ([]models.Tag, error)
}

type movieRepo struct{ db *gorm.DB }
type linkRepo struct{ db *gorm.DB }
type ratingRepo struct{ db *gorm.DB }
type tagRepo struct{ db *gorm.DB }

func NewMovieRepository(db *gorm.DB) MovieRepository   { return &movieRepo{db} }
func NewLinkRepository(db *gorm.DB) LinkRepository     { return &linkRepo{db} }
func NewRatingRepository(db *gorm.DB) RatingRepository { return &ratingRepo{db} }
func NewTagRepository(db *gorm.DB) TagRepository       { return &tagRepo{db} }

func (r *movieRepo) List(ctx context.Context, offset, limit int) ([]models.MovieListItem, error) {
	var items []models.MovieListItem
	err := r.db.WithContext(ctx).Order("movie_id").Offset(offset).Limit(limit).Find(&items).Error
	return items, err
}

func (r *movieRepo) GetByID(ctx context.Context, id int) (*models.Movie, error) {
	var m models.Movie
	if err := r.db.WithContext(ctx).First(&m, "movie_id = ?", id).Error; err != nil {
		return nil, mapErr(err)
	}
	return &m, nil
}

func (r *linkRepo) GetByMovieID(ctx context.Context, movieID int) (*models.Link, error) {
	var l models.Link
	if err := r.db.WithContext(ctx).First(&l, "movie_id = ?", movieID).Error; err != nil {
		return nil, mapErr(err)
	}
	return &l, nil
}

func (r *ratingRepo) ListByMovieID(ctx context.Context, movieID, offset, limit int) ([]models.Rating, error) {
	var out []models.Rating
	err := r.db.WithContext(ctx).Where("movie_id = ?", movieID).
		Order("user_id").Offset(offset).Limit(limit).Find(&out).Error
	return out, err
}

func (r *tagRepo) ListByMovieID(ctx context.Context, movieID int) ([]models.Tag, error) {
	var out []models.Tag
	err := r.db.WithContext(ctx).Where("movie_id = ?", movieID).Order("tag_id").Find(&out).Error
	return out, err
}

func mapErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
