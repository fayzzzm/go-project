-- =============================================================================
-- TEAMS DOMAIN
-- =============================================================================

CREATE SCHEMA IF NOT EXISTS teams;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- TABLE: teams.team
CREATE TABLE IF NOT EXISTS teams.team (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    status TEXT DEFAULT 'active',
    tenant_id TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT DEFAULT '',
    updated_by TEXT DEFAULT ''
);

-- INDEX: Active Teams
CREATE INDEX IF NOT EXISTS idx_teams_active ON teams.team(status) WHERE status = 'active';

-- TYPE: Request/Response DTOs
CREATE TYPE teams.team_request AS (
    id            UUID,
    name          TEXT,
    status        TEXT,
    tenant_id     TEXT,
    limit_val     INTEGER,
    offset_val    INTEGER
);

CREATE TYPE teams.team_response AS (
    id            UUID,
    name          TEXT,
    status        TEXT,
    tenant_id     TEXT,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ,
    created_by    TEXT,
    updated_by    TEXT
);

-- FUNCTIONS

-- Create
CREATE OR REPLACE FUNCTION teams.create(r teams.team_request)
RETURNS SETOF teams.team_response AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO teams.team (
        name, status, tenant_id
    )
    VALUES (
        TRIM(r.name), COALESCE(r.status, 'active'), r.tenant_id
    )
    RETURNING id INTO v_id;
    
    RETURN QUERY SELECT id, name, status, tenant_id, created_at, updated_at, created_by, updated_by
    FROM teams.team WHERE id = v_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Get By ID
CREATE OR REPLACE FUNCTION teams.get_by_id(r teams.team_request)
RETURNS SETOF teams.team_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, status, tenant_id, created_at, updated_at, created_by, updated_by
    FROM teams.team
    WHERE id = r.id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- List
CREATE OR REPLACE FUNCTION teams.list(r teams.team_request)
RETURNS SETOF teams.team_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, status, tenant_id, created_at, updated_at, created_by, updated_by
    FROM teams.team
    WHERE (r.tenant_id IS NULL OR tenant_id = r.tenant_id)
    ORDER BY created_at DESC
    LIMIT COALESCE(r.limit_val, 100)
    OFFSET COALESCE(r.offset_val, 0);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- Update
CREATE OR REPLACE FUNCTION teams.update(r teams.team_request)
RETURNS SETOF teams.team_response AS $$
BEGIN
    RETURN QUERY
    UPDATE teams.team
    SET 
        name = COALESCE(r.name, name),
        status = COALESCE(r.status, status),
        updated_at = NOW()
    WHERE id = r.id
    RETURNING id, name, status, tenant_id, created_at, updated_at, created_by, updated_by;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Delete
CREATE OR REPLACE FUNCTION teams.delete(r teams.team_request)
RETURNS VOID AS $$
BEGIN
    DELETE FROM teams.team WHERE id = r.id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

