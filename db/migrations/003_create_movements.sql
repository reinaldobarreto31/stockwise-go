-- Migration: 003_create_movements
-- Creates the stock movements table

CREATE TYPE movement_type AS ENUM ('in', 'out');

CREATE TABLE IF NOT EXISTS movements (
    id         BIGSERIAL      PRIMARY KEY,
    product_id BIGINT         NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    type       movement_type  NOT NULL,
    quantity   INTEGER        NOT NULL CHECK (quantity > 0),
    note       TEXT           NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_movements_product_id ON movements(product_id);
CREATE INDEX IF NOT EXISTS idx_movements_created_at ON movements(created_at DESC);
