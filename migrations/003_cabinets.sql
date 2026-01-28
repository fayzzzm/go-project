-- =============================================================================
-- CABINETS DOMAIN
-- Depends on: teams
-- =============================================================================

CREATE SCHEMA IF NOT EXISTS cabinets;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- TABLE: cabinets.cabinet
-- TABLE: cabinets.cabinet
CREATE TABLE IF NOT EXISTS cabinets.cabinet (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    location TEXT,
    machine_id TEXT,
    status TEXT DEFAULT 'available',
    team_id UUID REFERENCES teams.team(id),
    tenant_id UUID REFERENCES tenants.tenant(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID,
    updated_by UUID
);

-- INDEX: Available Cabinets
CREATE INDEX IF NOT EXISTS idx_cabinets_available ON cabinets.cabinet(status) WHERE status = 'available';

-- CABINET TYPES
CREATE TYPE cabinets.cabinet_request AS (
    id            UUID,
    name          TEXT,
    description   TEXT,
    location      TEXT,
    machine_id    TEXT,
    status        TEXT,
    team_id       UUID,
    tenant_id     UUID,
    limit_val     INTEGER,
    offset_val    INTEGER,
    user_id       UUID,
    created_by    UUID,
    updated_by    UUID
);

CREATE TYPE cabinets.cabinet_response AS (
    id          UUID,
    name        TEXT,
    description TEXT,
    location    TEXT,
    machine_id  TEXT,
    status      TEXT,
    team_id     UUID,
    tenant_id   UUID,
    created_at  TIMESTAMPTZ,
    created_by  UUID,
    updated_at  TIMESTAMPTZ,
    updated_by  UUID
);

-- CABINET FUNCTIONS
CREATE OR REPLACE FUNCTION cabinets.create(r cabinets.cabinet_request)
RETURNS SETOF cabinets.cabinet_response AS $$
DECLARE
    v_id UUID;
BEGIN
    -- Check Membership
    IF r.team_id IS NOT NULL THEN
        IF r.user_id IS NULL THEN
            RAISE EXCEPTION 'User ID is required when creating a cabinet for a team';
        END IF;
        
        IF NOT EXISTS (SELECT 1 FROM teams.members WHERE team_id = r.team_id AND user_id = r.user_id) THEN
             RAISE EXCEPTION 'User % is not a member of team %', r.user_id, r.team_id USING ERRCODE = 'P0001';
        END IF;
    END IF;

    INSERT INTO cabinets.cabinet (
        name, description, location, machine_id, team_id, tenant_id, created_by, updated_by
    )
    VALUES (
        r.name, r.description, r.location, r.machine_id, r.team_id, r.tenant_id, r.created_by, r.updated_by
    )
    RETURNING id INTO v_id;
    
    RETURN QUERY 
    SELECT id, name, description, location, machine_id, status, team_id, tenant_id, created_at, created_by, updated_at, updated_by
    FROM cabinets.cabinet WHERE id = v_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE OR REPLACE FUNCTION cabinets.get_by_id(r cabinets.cabinet_request)
RETURNS SETOF cabinets.cabinet_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, description, location, machine_id, status, team_id, tenant_id, created_at, created_by, updated_at, updated_by
    FROM cabinets.cabinet
    WHERE id = r.id AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

CREATE OR REPLACE FUNCTION cabinets.list(r cabinets.cabinet_request)
RETURNS SETOF cabinets.cabinet_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, description, location, machine_id, status, team_id, tenant_id, created_at, created_by, updated_at, updated_by
    FROM cabinets.cabinet
    WHERE (r.tenant_id IS NULL OR tenant_id = r.tenant_id)
    ORDER BY created_at DESC
    LIMIT COALESCE(r.limit_val, 100)
    OFFSET COALESCE(r.offset_val, 0);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

CREATE OR REPLACE FUNCTION cabinets.update(r cabinets.cabinet_request)
RETURNS SETOF cabinets.cabinet_response AS $$
BEGIN
    RETURN QUERY
    UPDATE cabinets.cabinet
    SET 
        name = COALESCE(r.name, name),
        description = COALESCE(r.description, description),
        location = COALESCE(r.location, location),
        machine_id = COALESCE(r.machine_id, machine_id),
        status = COALESCE(r.status, status),
        team_id = COALESCE(r.team_id, team_id),
        updated_at = NOW(),
        updated_by = COALESCE(r.updated_by, updated_by)
    WHERE id = r.id AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id)
    RETURNING id, name, description, location, machine_id, status, team_id, tenant_id, created_at, created_by, updated_at, updated_by;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE OR REPLACE FUNCTION cabinets.delete(r cabinets.cabinet_request)
RETURNS VOID AS $$
BEGIN
    DELETE FROM cabinets.cabinet WHERE id = r.id AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
