-- +goose Up
CREATE TABLE train_stations (
    id uuid NOT NULL,
    station_name VARCHAR 100 NOT NULL,
    station_code VARCHAR 50 NOT NULL
);


-- +goose Down
DROP TABLE train_stations;