# Edge Case Analysis & Security Requirements

This document tracks identified edge cases and security vulnerabilities that need to be addressed in the itsware API implementation.

## Authentication & Identity
### Login Timing Attack
- **Problem**: API returns 404 immediately if email doesn't exist, allowing email enumeration.
- **Requirement**: Always perform a dummy password hash check if the user is missing to maintain consistent response times.
- **Status**: Needs Implementation.

### User Creation Race Condition
- **Problem**: Duplicate emails might return raw 500 errors.
- **Requirement**: Catch DB unique constraint violations and return a 409 Conflict with a clear message.
- **Status**: Partially Handled (DB enforces uniqueness, but API response needs improvement).

## Teams & Membership
### Cross-Tenant Membership Injection
- **Problem**: A user could technically be added to a team belonging to a different tenant.
- **Requirement**: TeamUseCase.AddMember must verify that the UserID and TeamID both belong to the same TenantID.
- **Status**: Needs Implementation.

### Membership Privacy
- **Problem**: Users might be able to see who belongs to teams in other tenants.
- **Requirement**: Filter all team lookups strictly by the requester's TenantID.
- **Status**: Partially Handled (List functions have tenant filters, but validation checks are missing).

## Cabinets
### Orphaned Devices on Delete
- **Problem**: Deleting a cabinet might leave devices with invalid cabinet_id references.
- **Requirement**: 
    - **Option A**: Prevent deletion if devices are present (400 Bad Request).
    - **Option B**: Automatically unassign devices (set cabinet_id to NULL) in a transaction.
- **Status**: Needs Implementation.

### Global Machine ID Uniqueness
- **Problem**: Physical hardware (machine_id) should be unique globally across the system.
- **Requirement**: Add a unique constraint or validation to prevent two tenants from registering the same physical cabinet.
- **Status**: Needs Implementation.

## Devices
### Tenant-Profile Mismatch
- **Problem**: A device in Tenant A could be assigned to a Device Profile owned by Tenant B.
- **Requirement**: DeviceUseCase must validate that device_profile_id, cabinet_id, and team_id all share the same TenantID as the device itself.
- **Status**: Partially Handled (Basic tenant checks exist in Update/Delete, but deep assignment validation is missing).

### EPC (RFID) "Ghost" Lookups
- **Problem**: Hardware scanners often query EPCs that aren't yet in the database.
- **Requirement**: Ensure GetByEPC is optimized and returns a specific, predictable "Not Found" state that doesn't trigger error alerts in the hardware logs.
- **Status**: Partially Handled (Basic implementation exists).

## Global Implementation
### Deep Offset Performance
- **Problem**: High pagination offsets (e.g. 100,000+) degrade Postgres performance.
- **Requirement**: Maximize limit (currently 1000) and implement a max_offset cap or transition to cursor-based pagination for large tables.
- **Status**: Partially Handled (Limit is capped, but offset is not).

### Soft Deletion Strategy
- **Requirement**: Evaluate if Entities (Teams, Devices) should be "Soft Deleted" (flagged) instead of physically removed to preserve audit trails.
- **Status**: Needs Implementation.
