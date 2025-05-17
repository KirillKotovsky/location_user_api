package model

import (
	"time"

	locationpb "github.com/KirillKotovsky/location_proto/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DTO —  http request
type LocationUpdateRequest struct {
	Username  string    `json:"username" binding:"required,alphanum,min=4,max=16"`
	Latitude  float64   `json:"latitude" binding:"required"`
	Longitude float64   `json:"longitude" binding:"required"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// GORM-model
type UserLocation struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"index;not null"`
	Latitude  float64
	Longitude float64
	Timestamp time.Time `gorm:"not null"`
}

// DTO → GORM
func (r LocationUpdateRequest) ToEntity() UserLocation {
	t := r.Timestamp
	if t.IsZero() {
		t = time.Now().UTC()
	}
	return UserLocation{
		Username:  r.Username,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
		Timestamp: t,
	}
}

// GORM → DTO
func (u UserLocation) ToDTO() LocationUpdateRequest {
	return LocationUpdateRequest{
		Username:  u.Username,
		Latitude:  u.Latitude,
		Longitude: u.Longitude,
		Timestamp: u.Timestamp,
	}
}

func (l UserLocation) ToProto() *locationpb.Location {
	return &locationpb.Location{
		Username:  l.Username,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Timestamp: timestamppb.New(l.Timestamp),
	}
}
