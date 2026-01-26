-- =============================================================================
-- IDENTITY DOMAIN (Users)
-- =============================================================================

CREATE SCHEMA IF NOT EXISTS users;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- TABLE: users.user (Singular)
-- Stores all identity information
CREATE TABLE IF NOT EXISTS users.user (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    address TEXT,
    role TEXT DEFAULT 'user',
    phone TEXT,
    app_metadata JSONB DEFAULT '{}',
    user_metadata JSONB DEFAULT '{}',
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
    limit_val     INTEGER,
    offset_val    INTEGER
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
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ
);

-- FUNCTIONS

-- Create
CREATE OR REPLACE FUNCTION users.create(r users.user_request)
RETURNS SETOF users.user_response AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO users.user (
        email, name, address, phone, role, app_metadata, user_metadata
    )
    VALUES (
        LOWER(TRIM(r.email)), TRIM(r.name), TRIM(r.address), 
        r.phone, COALESCE(r.role, 'user'), 
        COALESCE(r.app_metadata, '{}'::jsonb), 
        COALESCE(r.user_metadata, '{}'::jsonb)
    )
    RETURNING id INTO v_id;
    
    RETURN QUERY SELECT id, email, name, address, phone, role, app_metadata, user_metadata, created_at, updated_at
    FROM users.user WHERE id = v_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Get By ID
CREATE OR REPLACE FUNCTION users.get_by_id(r users.user_request)
RETURNS SETOF users.user_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, email, name, address, phone, role, app_metadata, user_metadata, created_at, updated_at
    FROM users.user
    WHERE id = r.id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- List
CREATE OR REPLACE FUNCTION users.list(r users.user_request)
RETURNS SETOF users.user_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, email, name, address, phone, role, app_metadata, user_metadata, created_at, updated_at
    FROM users.user
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
    RETURNING id, email, name, address, phone, role, app_metadata, user_metadata, created_at, updated_at;
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
    SELECT id, email, name, address, phone, role, app_metadata, user_metadata, created_at, updated_at
    FROM users.user
    WHERE email = LOWER(TRIM(r.email));
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;
