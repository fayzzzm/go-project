-- =============================================================================
-- SEED DATA
-- =============================================================================

-- 1. Tenant
INSERT INTO tenants.tenant (id, name)
VALUES 
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'ItsWare'),
    ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Globex Inc')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

-- 2. Users (Admin and Regular)
INSERT INTO users.user (id, email, name, role, tenant_id, password)
VALUES 
    -- Admin for ItsWare (Password: password123)
    ('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'admin@itsware.com', 'ItsWare Admin', 'admin', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '$2a$10$F.K57ZES07MXSoenLRpSKuQc2SjUMQUGFbWhFdlN/Au8/0oe.NP8W'),
    -- User for ItsWare
    ('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'user@itsware.com', 'ItsWare User', 'user', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '$2a$10$F.K57ZES07MXSoenLRpSKuQc2SjUMQUGFbWhFdlN/Au8/0oe.NP8W'),
    -- Admin for Globex
    ('cccc3333-cccc-cccc-cccc-cccccccccccc', 'admin@globex.com', 'Globex Admin', 'admin', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', '$2a$10$F.K57ZES07MXSoenLRpSKuQc2SjUMQUGFbWhFdlN/Au8/0oe.NP8W')
ON CONFLICT (email) DO NOTHING;

-- 3. Teams (Only for Acme for now)
INSERT INTO teams.team (id, name, status, tenant_id, created_by)
VALUES 
    ('34c5a214-9e6b-4bb1-9f91-3328b222805c', 'Engineering', 'active', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa'),
    ('44c5a214-9e6b-4bb1-9f91-3328b222805d', 'Sales', 'active', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa')
ON CONFLICT DO NOTHING;

-- 4. Team Memberships
INSERT INTO teams.members (team_id, user_id, role)
VALUES 
    ('34c5a214-9e6b-4bb1-9f91-3328b222805c', 'bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'MEMBER'), -- Acme User -> Engineering
    ('34c5a214-9e6b-4bb1-9f91-3328b222805c', 'aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'OWNER')   -- Acme Admin -> Engineering
ON CONFLICT DO NOTHING;

-- 5. Device Profiles
INSERT INTO device_profiles.device_profile (id, name, description, tenant_id, created_by)
VALUES 
    ('55c5a214-9e6b-4bb1-9f91-3328b222805e', 'RFID Reader X1', 'Standard RFID Reader', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa')
ON CONFLICT DO NOTHING;

-- 6. Cabinets
INSERT INTO cabinets.cabinet (id, name, location, status, team_id, tenant_id, created_by)
VALUES 
    ('66c5a214-9e6b-4bb1-9f91-3328b222805f', 'Server Rack A', 'Data Center 1', 'available', '34c5a214-9e6b-4bb1-9f91-3328b222805c', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa')
ON CONFLICT DO NOTHING;

-- 7. Devices
INSERT INTO devices.device (id, name, serial_number, epc, device_profile_id, cabinet_id, team_id, tenant_id, created_by)
VALUES 
    ('77c5a214-9e6b-4bb1-9f91-3328b2228050', 'Reader-001', 'SN-001', 'EPC-001', '55c5a214-9e6b-4bb1-9f91-3328b222805e', '66c5a214-9e6b-4bb1-9f91-3328b222805f', '34c5a214-9e6b-4bb1-9f91-3328b222805c', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa')
ON CONFLICT (serial_number) DO NOTHING;
