-- =============================================================================
-- DEVICE PROFILES DOMAIN
-- =============================================================================

CREATE SCHEMA IF NOT EXISTS device_profiles;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- TABLE: device_profiles.device_profile
CREATE TABLE IF NOT EXISTS device_profiles.device_profile (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    tenant_id UUID REFERENCES tenants.tenant(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- TYPE: Request/Response DTOs
CREATE TYPE device_profiles.device_profile_request AS (
    id            UUID,
    name          TEXT,
    description   TEXT,
    tenant_id     UUID,
    limit_val     INTEGER,
    offset_val    INTEGER
);

CREATE TYPE device_profiles.device_profile_response AS (
    id            UUID,
    name          TEXT,
    description   TEXT,
    tenant_id     UUID,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ
);

-- FUNCTIONS

-- Create
CREATE OR REPLACE FUNCTION device_profiles.create(r device_profiles.device_profile_request)
RETURNS SETOF device_profiles.device_profile_response AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO device_profiles.device_profile (
        name, description, tenant_id
    )
    VALUES (
        TRIM(r.name), TRIM(r.description), r.tenant_id
    )
    RETURNING id INTO v_id;
    
    RETURN QUERY SELECT id, name, description, tenant_id, created_at, updated_at
    FROM device_profiles.device_profile WHERE id = v_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Get By ID
CREATE OR REPLACE FUNCTION device_profiles.get_by_id(r device_profiles.device_profile_request)
RETURNS SETOF device_profiles.device_profile_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, description, tenant_id, created_at, updated_at
    FROM device_profiles.device_profile
    WHERE id = r.id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- List
CREATE OR REPLACE FUNCTION device_profiles.list(r device_profiles.device_profile_request)
RETURNS SETOF device_profiles.device_profile_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, description, tenant_id, created_at, updated_at
    FROM device_profiles.device_profile
    WHERE (r.tenant_id IS NULL OR tenant_id = r.tenant_id)
    ORDER BY created_at DESC
    LIMIT COALESCE(r.limit_val, 100)
    OFFSET COALESCE(r.offset_val, 0);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- Update
CREATE OR REPLACE FUNCTION device_profiles.update(r device_profiles.device_profile_request)
RETURNS SETOF device_profiles.device_profile_response AS $$
BEGIN
    RETURN QUERY
    UPDATE device_profiles.device_profile
    SET 
        name = COALESCE(r.name, name),
        description = COALESCE(r.description, description),
        updated_at = NOW()
    WHERE id = r.id
    RETURNING id, name, description, tenant_id, created_at, updated_at;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Delete
CREATE OR REPLACE FUNCTION device_profiles.delete(r device_profiles.device_profile_request)
RETURNS VOID AS $$
BEGIN
    DELETE FROM device_profiles.device_profile WHERE id = r.id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

