-- =============================================================================
-- IDENTITY DOMAIN (Users)
-- Usage: Depends on 000_tenants.sql
-- =============================================================================

-- USERS
CREATE SCHEMA IF NOT EXISTS users;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users.user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email CITEXT UNIQUE NOT NULL,
    name TEXT,
    address TEXT,
    role TEXT DEFAULT 'user',
    phone TEXT,
    app_metadata JSONB DEFAULT '{}',
    user_metadata JSONB DEFAULT '{}',
    password TEXT,
    tenant_id UUID REFERENCES tenants.tenant(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- TYPE: Request/Response DTOs
CREATE TYPE users.user_request AS (
    id            UUID,
    email         TEXT,
    name          TEXT,
    address       TEXT,
    phone         TEXT,
    role          TEXT,
    app_metadata  JSONB,
    user_metadata JSONB,
    password      TEXT,
    limit_val     INTEGER,
    offset_val    INTEGER,
    team_id       UUID, -- Added for future use
    tenant_id     UUID  -- Changed to UUID
);

CREATE TYPE users.user_response AS (
    id            UUID,
    email         TEXT,
    name          TEXT,
    address       TEXT,
    phone         TEXT,
    role          TEXT,
    app_metadata  JSONB,
    user_metadata JSONB,
    password      TEXT, -- Restored for struct compatibility, returns empty
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    tenant_id     UUID
);

-- FUNCTIONS (Updated to use tenant_id logic)
-- Create
CREATE OR REPLACE FUNCTION users.create(r users.user_request)
RETURNS SETOF users.user_response AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO users.user (
        email, name, address, phone, role, app_metadata, user_metadata, password, tenant_id
    )
    VALUES (
        TRIM(r.email), TRIM(r.name), TRIM(r.address), 
        r.phone, COALESCE(r.role, 'user'), 
        COALESCE(r.app_metadata, '{}'::jsonb), 
        COALESCE(r.user_metadata, '{}'::jsonb),
        r.password,
        r.tenant_id
    )
    RETURNING id INTO v_id;
    
    RETURN QUERY SELECT id, email, name, address, phone, role, app_metadata, user_metadata, ''::text as password, created_at, updated_at, tenant_id
    FROM users.user WHERE id = v_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Get By ID
CREATE OR REPLACE FUNCTION users.get_by_id(r users.user_request)
RETURNS SETOF users.user_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, email, name, address, phone, role, app_metadata, user_metadata, ''::text as password, created_at, updated_at, tenant_id
    FROM users.user
    WHERE id = r.id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- List (Tenant filter only)
CREATE OR REPLACE FUNCTION users.list(r users.user_request)
RETURNS SETOF users.user_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, email, name, address, phone, role, app_metadata, user_metadata, ''::text as password, created_at, updated_at, tenant_id
    FROM users.user
    WHERE (r.tenant_id IS NULL OR tenant_id = r.tenant_id)
    ORDER BY created_at DESC
    LIMIT COALESCE(r.limit_val, 100)
    OFFSET COALESCE(r.offset_val, 0);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- Update
CREATE OR REPLACE FUNCTION users.update(r users.user_request)
RETURNS SETOF users.user_response AS $$
BEGIN
    RETURN QUERY
    UPDATE users.user
    SET 
        name = COALESCE(r.name, name),
        address = COALESCE(r.address, address),
        phone = COALESCE(r.phone, phone),
        role = COALESCE(r.role, role),
        app_metadata = COALESCE(r.app_metadata, app_metadata),
        user_metadata = COALESCE(r.user_metadata, user_metadata),
        updated_at = NOW()
    WHERE id = r.id
    RETURNING id, email, name, address, phone, role, app_metadata, user_metadata, ''::text as password, created_at, updated_at, tenant_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Delete
CREATE OR REPLACE FUNCTION users.delete(r users.user_request)
RETURNS VOID AS $$
BEGIN
    DELETE FROM users.user WHERE id = r.id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Get By Email
CREATE OR REPLACE FUNCTION users.get_by_email(r users.user_request)
RETURNS SETOF users.user_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, email, name, address, phone, role, app_metadata, user_metadata, ''::text as password, created_at, updated_at, tenant_id
    FROM users.user
    WHERE email = TRIM(r.email);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- Get For Login
CREATE OR REPLACE FUNCTION users.get_for_login(p_email TEXT)
RETURNS TABLE (
    id UUID, email TEXT, name TEXT, address TEXT, phone TEXT,
    role TEXT, app_metadata JSONB, user_metadata JSONB,
    password TEXT, created_at TIMESTAMPTZ, updated_at TIMESTAMPTZ,
    tenant_id UUID
) AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.email, u.name, u.address, u.phone,
           u.role, u.app_metadata, u.user_metadata,
           u.password,
           u.created_at, u.updated_at, u.tenant_id
    FROM users.user u
    WHERE u.email = TRIM(p_email);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- Check Tenant Access
CREATE OR REPLACE FUNCTION users.check_tenant_access(p_user_id UUID, p_tenant_id UUID)
RETURNS BOOLEAN AS $$
DECLARE
    v_exists BOOLEAN;
BEGIN
    IF p_user_id IS NULL OR p_tenant_id IS NULL THEN
        RETURN FALSE;
    END IF;
    SELECT EXISTS(SELECT 1 FROM users.user WHERE id = p_user_id AND tenant_id = p_tenant_id) INTO v_exists;
    RETURN v_exists;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;
