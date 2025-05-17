package model

type NearbyUsersRequest struct {
	Latitude  float64 `form:"latitude" binding:"required"`  // alias for lat
	Lat       float64 `form:"lat" binding:"-"`              // optional alias
	Longitude float64 `form:"longitude" binding:"required"` // alias for lon
	Lon       float64 `form:"lon" binding:"-"`              // optional alias
	Radius    float64 `form:"radius" binding:"required"`    // in kilometers
	Page      int     `form:"page,default=1"`
	Limit     int     `form:"limit,default=10"`
	Offset    int     `form:"offset,default=0"`
}
