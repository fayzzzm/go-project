-- =============================================================================
-- TENANT DOMAIN
-- =============================================================================

CREATE SCHEMA IF NOT EXISTS tenants;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tenants.tenant (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Seed Default Tenant
INSERT INTO tenants.tenant (id, name) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Default Tenant') ON CONFLICT DO NOTHING;
