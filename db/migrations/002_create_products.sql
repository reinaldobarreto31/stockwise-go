-- Migration: 002_create_products
-- Creates the products table for stock items

CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    sku         VARCHAR(100)   NOT NULL UNIQUE,
    category    VARCHAR(100)   NOT NULL DEFAULT '',
    description TEXT           NOT NULL DEFAULT '',
    price       NUMERIC(12, 2) NOT NULL DEFAULT 0,
    stock       INTEGER        NOT NULL DEFAULT 0,
    min_stock   INTEGER        NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_sku      ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
