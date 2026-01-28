-- =============================================================================
-- TENANT DOMAIN
-- =============================================================================

CREATE SCHEMA IF NOT EXISTS tenants;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS tenants.tenant (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- TYPE: Request/Response DTOs
CREATE TYPE tenants.tenant_request AS (
    id UUID,
    name TEXT,
    limit_val INTEGER,
    offset_val INTEGER
);

CREATE TYPE tenants.tenant_response AS (
    id UUID,
    name TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

-- FUNCTIONS
-- List
CREATE OR REPLACE FUNCTION tenants.list(r tenants.tenant_request)
RETURNS SETOF tenants.tenant_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, created_at, updated_at
    FROM tenants.tenant
    ORDER BY created_at DESC
    LIMIT COALESCE(r.limit_val, 100)
    OFFSET COALESCE(r.offset_val, 0);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- Seed Default Tenant
INSERT INTO tenants.tenant (id, name) VALUES ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Default Tenant') ON CONFLICT DO NOTHING;
