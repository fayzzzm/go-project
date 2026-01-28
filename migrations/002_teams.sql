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
    tenant_id UUID REFERENCES tenants.tenant(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by TEXT DEFAULT '',
    updated_by TEXT DEFAULT ''
);

-- INDEX: Active Teams
CREATE INDEX IF NOT EXISTS idx_teams_active ON teams.team(status) WHERE status = 'active';

-- MEMBERSHIP TABLE (Replaces old user_teams)
CREATE TABLE IF NOT EXISTS teams.members (
    team_id UUID NOT NULL REFERENCES teams.team(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users.user(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'MEMBER',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (team_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_team_members_user ON teams.members(user_id);

-- TYPE: Request/Response DTOs
CREATE TYPE teams.team_request AS (
    id            UUID,
    name          TEXT,
    status        TEXT,
    tenant_id     UUID,
    limit_val     INTEGER,
    offset_val    INTEGER,
    user_id       UUID
);

CREATE TYPE teams.team_response AS (
    id            UUID,
    name          TEXT,
    status        TEXT,
    tenant_id     UUID,
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
    WHERE id = r.id AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- List
CREATE OR REPLACE FUNCTION teams.list(r teams.team_request)
RETURNS SETOF teams.team_response AS $$
BEGIN
    -- Check Tenant Access
    IF r.user_id IS NOT NULL AND r.tenant_id IS NOT NULL THEN
        IF NOT users.check_tenant_access(r.user_id, r.tenant_id) THEN
             RAISE EXCEPTION 'User % does not belong to tenant %', r.user_id, r.tenant_id USING ERRCODE = 'P0001';
        END IF;
    END IF;

    RETURN QUERY
    SELECT id, name, status, tenant_id, created_at, updated_at, created_by, updated_by
    FROM teams.team
    WHERE tenant_id = r.tenant_id
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
    WHERE id = r.id AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id)
    RETURNING id, name, status, tenant_id, created_at, updated_at, created_by, updated_by;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Delete
CREATE OR REPLACE FUNCTION teams.delete(r teams.team_request)
RETURNS VOID AS $$
BEGIN
    DELETE FROM teams.team WHERE id = r.id AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- MEMBER TYPES
CREATE TYPE teams.member_request AS (
    team_id    UUID,
    user_id    UUID,
    role       TEXT
);

CREATE TYPE teams.member_response AS (
    team_id    UUID,
    user_id    UUID,
    role       TEXT,
    created_at TIMESTAMPTZ
);

-- Add Member
CREATE OR REPLACE FUNCTION teams.add_member(r teams.member_request)
RETURNS SETOF teams.member_response AS $$
BEGIN
    RETURN QUERY
    INSERT INTO teams.members (team_id, user_id, role)
    VALUES (r.team_id, r.user_id, COALESCE(r.role, 'MEMBER'))
    RETURNING team_id, user_id, role::TEXT, created_at;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Get User Teams
CREATE OR REPLACE FUNCTION teams.get_for_user(p_user_id UUID)
RETURNS SETOF teams.team_response AS $$
BEGIN
    RETURN QUERY
    SELECT t.id, t.name, t.status, t.tenant_id, t.created_at, t.updated_at, t.created_by, t.updated_by
    FROM teams.team t
    JOIN teams.members m ON m.team_id = t.id
    WHERE m.user_id = p_user_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

-- Update users.list to support filtering by team_id
-- We consolidate it here because it depends on the teams.members table
CREATE OR REPLACE FUNCTION users.list(r users.user_request)
RETURNS SETOF users.user_response AS $$
BEGIN
    RETURN QUERY
    SELECT u.id, u.email, u.name, u.address, u.phone, u.role, u.app_metadata, u.user_metadata, ''::text as password, u.created_at, u.updated_at, u.tenant_id
    FROM users.user u
    WHERE u.tenant_id = r.tenant_id
      AND (r.team_id IS NULL OR EXISTS (
          SELECT 1 FROM teams.members m WHERE m.user_id = u.id AND m.team_id = r.team_id
      ))
    ORDER BY u.created_at DESC
    LIMIT COALESCE(r.limit_val, 100)
    OFFSET COALESCE(r.offset_val, 0);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

