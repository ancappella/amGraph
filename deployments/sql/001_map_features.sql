-- Map features: GeoJSON stored as JSONB (Point, LineString, Polygon, etc.)
CREATE TABLE IF NOT EXISTS map_features (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    geojson JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_map_features_name ON map_features (name);
CREATE INDEX IF NOT EXISTS idx_map_features_geojson ON map_features USING GIN (geojson);
