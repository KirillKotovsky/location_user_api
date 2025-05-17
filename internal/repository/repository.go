package repository

import (
	"context"

	"github.com/KirillKotovsky/location_user_api/internal/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	SaveLocation(ctx context.Context, loc model.UserLocation) error
	InitSchema(ctx context.Context) error
	FindNearby(ctx context.Context, req model.NearbyUsersRequest) ([]model.UserLocation, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r *repo) InitSchema(ctx context.Context) error {
	const enableExtensions = `
		CREATE EXTENSION IF NOT EXISTS cube;
		CREATE EXTENSION IF NOT EXISTS earthdistance;
	`
	if _, err := r.db.ExecContext(ctx, enableExtensions); err != nil {
		return err
	}
	const tableCheck = `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' AND table_name = 'user_locations'
		);
	`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, tableCheck); err != nil {
		return err
	}

	if exists {
		return nil
	}

	createTable := `
		CREATE TABLE user_locations (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL
		);
		CREATE INDEX idx_username_time ON user_locations(username, timestamp);
	`

	_, err := r.db.ExecContext(ctx, createTable)
	return err
}

func (r *repo) SaveLocation(ctx context.Context, loc model.UserLocation) error {
	const insertQuery = `
		INSERT INTO user_locations (username, latitude, longitude, timestamp)
		VALUES ($1, $2, $3, $4);
	`
	_, err := r.db.ExecContext(ctx, insertQuery, loc.Username, loc.Latitude, loc.Longitude, loc.Timestamp)
	return err
}

func (r *repo) FindNearby(ctx context.Context, req model.NearbyUsersRequest) ([]model.UserLocation, error) {
	const query = `
		SELECT username, latitude, longitude, timestamp
		FROM user_locations
		WHERE earth_box(ll_to_earth($1, $2), $3) @> ll_to_earth(latitude, longitude)
		ORDER BY timestamp DESC
		OFFSET $4 LIMIT $5;
	`

	offset := (req.Page - 1) * req.Limit
	earthRadiusMeters := req.Radius * 1000 // convert km to meters

	var locations []model.UserLocation
	err := r.db.SelectContext(ctx, &locations, query,
		req.Latitude, req.Longitude, earthRadiusMeters, offset, req.Limit)

	return locations, err
}
