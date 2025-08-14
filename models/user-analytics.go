package models

import "time"

type Session struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
	VisitCount int       `json:"visit_count"`
	UserAgent  string    `json:"user_agent"`
	IP         string    `json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	// 一对多关系，但不用外键约束（性能考虑）
	VisitRecords []VisitRecord `gorm:"foreignKey:SessionID;constraint:OnDelete:SET NULL"`
}

type VisitRecord struct {
	ID            uint         `json:"id" gorm:"primaryKey"`
	GeoLocation   *GeoLocation `json:"geo_location" gorm:"foreignKey:GeoLocationId"` // Associated GeoLocation
	GeoLocationId *uint        `json:"geo_location_id"`                              // Foreign key to GeoLocation
	SessionID     uint         `json:"session_id"`                                   // Unique session identifier
	ShortURLID    string       `json:"short_url_id"`                                 // Unique short URL identifier
	VisitTime     time.Time    `json:"visit_time"`                                   // Time of visit
	CreatedAt     time.Time    `json:"create_at"`
	Session       *Session     `gorm:"foreignKey:SessionID;constraint:OnDelete:SET NULL"`
}

type GeoLocation struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Ip          string  `json:"ip"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Org         string  `json:"org"`
	Isp         string  `json:"isp"`
	As          string  `json:"as"`
	Timezone    string  `json:"timezone"`
}
