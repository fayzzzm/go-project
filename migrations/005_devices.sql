-- =============================================================================
-- DEVICES DOMAIN
-- Depends on: teams, cabinets, device_profiles
-- =============================================================================

CREATE SCHEMA IF NOT EXISTS devices;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- TABLE: devices.device
CREATE TABLE IF NOT EXISTS devices.device (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    serial_number TEXT UNIQUE NOT NULL,
    epc TEXT,
    device_profile_id UUID REFERENCES device_profiles.device_profile(id),
    cabinet_id UUID REFERENCES cabinets.cabinet(id),
    team_id UUID REFERENCES teams.team(id),
    tenant_id UUID REFERENCES tenants.tenant(id),
    status TEXT DEFAULT 'available',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- INDEX: Available Devices
CREATE INDEX IF NOT EXISTS idx_devices_available ON devices.device(status) WHERE status = 'available';

-- DEVICE TYPES
CREATE TYPE devices.device_request AS (
    id                 UUID,
    name               TEXT,
    description        TEXT,
    serial_number      TEXT,
    epc                TEXT,
    device_profile_id  UUID,
    cabinet_id         UUID,
    team_id            UUID,
    tenant_id          UUID,
    limit_val          INTEGER,
    offset_val         INTEGER
);

CREATE TYPE devices.device_response AS (
    id                 UUID,
    name               TEXT,
    description        TEXT,
    serial_number      TEXT,
    epc                TEXT,
    device_profile_id  UUID,
    cabinet_id         UUID,
    team_id            UUID,
    tenant_id          UUID,
    status             TEXT,
    created_at         TIMESTAMPTZ
);

-- DEVICE FUNCTIONS
CREATE OR REPLACE FUNCTION devices.create(r devices.device_request)
RETURNS SETOF devices.device_response AS $$
DECLARE
    v_id UUID;
BEGIN
    INSERT INTO devices.device (
        name, description, serial_number, epc, 
        device_profile_id, cabinet_id, team_id, tenant_id
    )
    VALUES (
        r.name, r.description, r.serial_number, r.epc, 
        r.device_profile_id, r.cabinet_id, r.team_id, r.tenant_id
    )
    RETURNING id INTO v_id;
    
    RETURN QUERY SELECT 
        id, name, description, serial_number, epc, 
        device_profile_id, cabinet_id, team_id, tenant_id, status, created_at 
    FROM devices.device WHERE id = v_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE OR REPLACE FUNCTION devices.get_by_id(r devices.device_request)
RETURNS SETOF devices.device_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, description, serial_number, epc, 
           device_profile_id, cabinet_id, team_id, tenant_id, status, created_at
    FROM devices.device
    WHERE id = r.id AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

CREATE OR REPLACE FUNCTION devices.list(r devices.device_request)
RETURNS SETOF devices.device_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, description, serial_number, epc, 
           device_profile_id, cabinet_id, team_id, tenant_id, status, created_at
    FROM devices.device
    WHERE (r.tenant_id IS NULL OR tenant_id = r.tenant_id)
    ORDER BY created_at DESC
    LIMIT COALESCE(r.limit_val, 100)
    OFFSET COALESCE(r.offset_val, 0);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

CREATE OR REPLACE FUNCTION devices.update(r devices.device_request)
RETURNS SETOF devices.device_response AS $$
BEGIN
    RETURN QUERY
    UPDATE devices.device
    SET 
        name = COALESCE(r.name, name),
        description = COALESCE(r.description, description),
        epc = COALESCE(r.epc, epc),
        device_profile_id = COALESCE(r.device_profile_id, device_profile_id),
        cabinet_id = COALESCE(r.cabinet_id, cabinet_id),
        team_id = COALESCE(r.team_id, team_id),
        updated_at = NOW()
    WHERE id = r.id AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id)
    RETURNING id, name, description, serial_number, epc, 
              device_profile_id, cabinet_id, team_id, tenant_id, status, created_at;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE OR REPLACE FUNCTION devices.delete(r devices.device_request)
RETURNS VOID AS $$
BEGIN
    DELETE FROM devices.device WHERE id = r.id AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE OR REPLACE FUNCTION devices.get_by_epc(r devices.device_request)
RETURNS SETOF devices.device_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, description, serial_number, epc, 
           device_profile_id, cabinet_id, team_id, tenant_id, status, created_at
    FROM devices.device
    WHERE epc = r.epc AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

CREATE OR REPLACE FUNCTION devices.get_by_serial_number(r devices.device_request)
RETURNS SETOF devices.device_response AS $$
BEGIN
    RETURN QUERY
    SELECT id, name, description, serial_number, epc, 
           device_profile_id, cabinet_id, team_id, tenant_id, status, created_at
    FROM devices.device
    WHERE serial_number = r.serial_number AND (r.tenant_id IS NULL OR tenant_id = r.tenant_id);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER STABLE;

CREATE OR REPLACE FUNCTION devices.bulk_create(r_list devices.device_request[])
RETURNS SETOF devices.device_response AS $$
BEGIN
    RETURN QUERY
    INSERT INTO devices.device (
        name, description, serial_number, epc, 
        device_profile_id, cabinet_id, team_id, tenant_id
    )
    SELECT 
        r.name, r.description, r.serial_number, r.epc, 
        r.device_profile_id, r.cabinet_id, r.team_id, r.tenant_id
    FROM unnest(r_list) r
    RETURNING id, name, description, serial_number, epc, 
              device_profile_id, cabinet_id, team_id, tenant_id, status, created_at;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;
