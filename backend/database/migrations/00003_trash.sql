-- +goose Up
CREATE TABLE trashes (
    trash_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    trash_name VARCHAR(100) NOT NULL,
    trash_code VARCHAR(50) NOT NULL,
    trash_json JSONB,
    trash_created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    trash_updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (trash_id)
);

-- +goose Down
DROP TABLE trashes;