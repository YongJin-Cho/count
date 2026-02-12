# Module Specification: auth-module

## 1. Overview
- **Role**: Handles authentication and authorization for the system, specifically validating Bearer Tokens for API access.
- **Build Output**: Library (shared package/middleware).

## 2. Providing Interfaces
- **Auth Interface**:
    - `ValidateToken(token string) (bool, error)`: Checks if the provided Bearer Token is valid.
    - `IsAuthorized(token string, resource string) (bool, error)`: Checks if the token has permissions for a specific resource/action.

## 3. Functional Requirements
- **User Story**: 
    - Related to FR-01-01 (API Security).
- **Core Logic**:
    - Parse "Authorization" header to extract Bearer Token.
    - Validate token authenticity (e.g., JWT signature check).
    - If token is missing or invalid, return error state corresponding to 401 Unauthorized (FR-01-01-02).
    - If token is valid but lacks permissions, return error state corresponding to 403 Forbidden (FR-01-01-03).

## 4. Dependencies
- **Reference Modules**: None.
- **Technologies Used**: Go, JWT.

## 5. Acceptance Tests
- [ ] Returns "Unauthorized" when no token is provided.
- [ ] Returns "Unauthorized" for expired or malformed tokens.
- [ ] Returns "Forbidden" for valid tokens without `collect` permission.
- [ ] Returns "Success" and allows access for valid tokens with correct permissions.
