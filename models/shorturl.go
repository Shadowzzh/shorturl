package models

import "time"

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

type ShortURL struct {
	ID          int64      `json:"id" gorm:"primaryKey"`
	Code        string     `json:"code" gorm:"uniqueIndex;size:10"`
	VisitCount  int        `json:"visit_count"`  // Number of times the shortened URL has been visited
	OriginalURL string     `json:"original_url"` // The original URL that was shortened
	CreatedAt   time.Time  `json:"created_at"`   // Timestamp when the shortened URL was created
	LastVisit   *time.Time `json:"last_visit"`   // Timestamp of the last visit to the shortened URL
}
