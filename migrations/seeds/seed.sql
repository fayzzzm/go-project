-- Seeding Data for Development
-- Usage: psql -d itsware -f migrations/seeds/seed.sql

-- 1. TEAMS
INSERT INTO teams.team (id, name, status, tenant_id) VALUES
    ('11111111-1111-1111-1111-111111111111', 'IT Operations', 'active', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
    ('22222222-2222-2222-2222-222222222222', 'Warehouse Logistics', 'active', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11')
ON CONFLICT (id) DO NOTHING;

-- 2. USERS
INSERT INTO users.user (id, email, name, role, tenant_id) VALUES
    ('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'admin@itsware.com', 'Super Admin', 'admin', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
    ('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'tech@itsware.com', 'Field Technician', 'user', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
    ('cccc3333-cccc-cccc-cccc-cccccccccccc', 'tenant_user@itsware.com', 'Tenant User', 'user', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11')
ON CONFLICT (email) DO NOTHING;

-- 3. CABINETS
INSERT INTO cabinets.cabinet (id, name, location, status, team_id, description) VALUES
    ('cccc3333-cccc-cccc-cccc-cccccccccccc', 'Server Rack Alpha', 'Data Center Row 1', 'available', '11111111-1111-1111-1111-111111111111', 'Main Storage'),
    ('dddd4444-dddd-dddd-dddd-dddddddddddd', 'Storage Unit B', 'Warehouse Loading Dock', 'maintenance', '22222222-2222-2222-2222-222222222222', 'Overflow Storage')
ON CONFLICT (id) DO NOTHING;

-- 4. DEVICE PROFILES
INSERT INTO device_profiles.device_profile (id, name, description) VALUES
    ('eeee5555-eeee-eeee-eeee-eeeeeeeeeeee', 'iPhone 15 Pro', 'Apple Smartphone 2023'),
    ('ffff6666-ffff-ffff-ffff-ffffffffffff', 'Zebra Scanner TC52', 'Rugged Barcode Scanner')
ON CONFLICT (id) DO NOTHING;

-- 5. DEVICES
INSERT INTO devices.device (name, serial_number, epc, device_profile_id, cabinet_id, team_id, status) VALUES
    ('iPhone 15 - Test Unit 1', 'SN-IPHONE-001', 'EPC-001', 'eeee5555-eeee-eeee-eeee-eeeeeeeeeeee', 'cccc3333-cccc-cccc-cccc-cccccccccccc', '11111111-1111-1111-1111-111111111111', 'available'),
    ('Zebra Scanner 001', 'SN-ZEBRA-001', 'EPC-002', 'ffff6666-ffff-ffff-ffff-ffffffffffff', 'dddd4444-dddd-dddd-dddd-dddddddddddd', '22222222-2222-2222-2222-222222222222', 'in_use')
ON CONFLICT (serial_number) DO NOTHING;
