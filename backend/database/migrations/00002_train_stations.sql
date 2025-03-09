-- +goose Up
CREATE TABLE train_stations (
    station_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    station_name VARCHAR(100) NOT NULL,
    station_code VARCHAR(50) NOT NULL,
    PRIMARY KEY (station_id),
    UNIQUE (station_name, station_code) -- Уникальный составной ключ
);

-- +goose Down
DROP TABLE train_stations;