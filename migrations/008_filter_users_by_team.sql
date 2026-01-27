-- Update users.list to support filtering by team_id
-- Dependent on 002_teams.sql (creates teams.members) and 001_users.sql (creates user type with team_id/tenant_id)

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
