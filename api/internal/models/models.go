package models

// Structure and attribute names match the database tables and columns.

type Movie struct {
	MovieID int    `gorm:"column:movie_id;primaryKey" json:"movie_id"`
	Title   string `gorm:"column:title" json:"title"`
	Genres  string `gorm:"column:genres" json:"genres"`
}

func (Movie) TableName() string { return "movies" }

// MovieListItem is the list projection of movies (ID and title only).
type MovieListItem struct {
	MovieID int    `gorm:"column:movie_id;primaryKey" json:"movie_id"`
	Title   string `gorm:"column:title" json:"title"`
}

func (MovieListItem) TableName() string { return "movies" }

type Link struct {
	MovieID int    `gorm:"column:movie_id;primaryKey" json:"movie_id"`
	ImdbID  string `gorm:"column:imdb_id" json:"imdb_id"`
	TmdbID  *int   `gorm:"column:tmdb_id" json:"tmdb_id"`
}

func (Link) TableName() string { return "links" }

type Rating struct {
	UserID  int     `gorm:"column:user_id;primaryKey" json:"user_id"`
	MovieID int     `gorm:"column:movie_id;primaryKey" json:"movie_id"`
	Rating  float64 `gorm:"column:rating" json:"rating"`
	RatedAt int64   `gorm:"column:rated_at" json:"rated_at"`
}

func (Rating) TableName() string { return "ratings" }

type Tag struct {
	TagID    int    `gorm:"column:tag_id;primaryKey" json:"tag_id"`
	UserID   int    `gorm:"column:user_id" json:"user_id"`
	MovieID  int    `gorm:"column:movie_id" json:"movie_id"`
	Tag      string `gorm:"column:tag" json:"tag"`
	TaggedAt int64  `gorm:"column:tagged_at" json:"tagged_at"`
}

func (Tag) TableName() string { return "tags" }
