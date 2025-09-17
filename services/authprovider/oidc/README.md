# Generic OpenID Connect (OIDC) Provider

This package provides a generic OpenID Connect authentication provider for the Go IAM system. It can work with any OIDC-compliant identity provider.

## Features

- **Generic OIDC Support**: Works with any OpenID Connect-compliant provider
- **OAuth2 Authorization Code Flow**: Full support for the standard OAuth2 flow
- **Refresh Token Support**: Automatic token refresh for long-lived sessions
- **UserInfo Integration**: Retrieves user profile information from the OIDC UserInfo endpoint
- **Flexible Configuration**: Configurable endpoints, scopes, and parameters

## Configuration Parameters

When creating an OIDC auth provider in the Go IAM system, you need to configure the following parameters:

### Required Parameters

- `@OIDC/CLIENT_ID`: OAuth2 client identifier from your OIDC provider
- `@OIDC/CLIENT_SECRET`: OAuth2 client secret from your OIDC provider
- `@OIDC/REDIRECT_URL`: The callback URL for your application
- `@OIDC/AUTHORIZATION_URL`: The authorization endpoint of your OIDC provider
- `@OIDC/TOKEN_URL`: The token endpoint of your OIDC provider
- `@OIDC/USERINFO_URL`: The UserInfo endpoint of your OIDC provider

## Usage Examples

### Example: Auth0 Configuration

```json
{
  "name": "Auth0 Provider",
  "provider": "OIDC",
  "params": [
    {
      "key": "@OIDC/CLIENT_ID",
      "value": "your-auth0-client-id"
    },
    {
      "key": "@OIDC/CLIENT_SECRET",
      "value": "your-auth0-client-secret"
    },
    {
      "key": "@OIDC/REDIRECT_URL",
      "value": "https://yourapp.com/auth/callback"
    },
    {
      "key": "@OIDC/AUTHORIZATION_URL",
      "value": "https://yourdomain.auth0.com/authorize"
    },
    {
      "key": "@OIDC/TOKEN_URL",
      "value": "https://yourdomain.auth0.com/oauth/token"
    },
    {
      "key": "@OIDC/USERINFO_URL",
      "value": "https://yourdomain.auth0.com/userinfo"
    },
    {
      "key": "@OIDC/SCOPES",
      "value": "openid profile email"
    }
  ]
}
```

### Example: Keycloak Configuration

```json
{
  "name": "Keycloak Provider",
  "provider": "OIDC",
  "params": [
    {
      "key": "@OIDC/CLIENT_ID",
      "value": "your-keycloak-client-id"
    },
    {
      "key": "@OIDC/CLIENT_SECRET",
      "value": "your-keycloak-client-secret"
    },
    {
      "key": "@OIDC/REDIRECT_URL",
      "value": "https://yourapp.com/auth/callback"
    },
    {
      "key": "@OIDC/AUTHORIZATION_URL",
      "value": "https://keycloak.example.com/auth/realms/yourrealm/protocol/openid-connect/auth"
    },
    {
      "key": "@OIDC/TOKEN_URL",
      "value": "https://keycloak.example.com/auth/realms/yourrealm/protocol/openid-connect/token"
    },
    {
      "key": "@OIDC/USERINFO_URL",
      "value": "https://keycloak.example.com/auth/realms/yourrealm/protocol/openid-connect/userinfo"
    }
  ]
}
```

### Example: Okta Configuration

```json
{
  "name": "Okta Provider",
  "provider": "OIDC",
  "params": [
    {
      "key": "@OIDC/CLIENT_ID",
      "value": "your-okta-client-id"
    },
    {
      "key": "@OIDC/CLIENT_SECRET",
      "value": "your-okta-client-secret"
    },
    {
      "key": "@OIDC/REDIRECT_URL",
      "value": "https://yourapp.com/auth/callback"
    },
    {
      "key": "@OIDC/AUTHORIZATION_URL",
      "value": "https://yourdomain.okta.com/oauth2/default/v1/authorize"
    },
    {
      "key": "@OIDC/TOKEN_URL",
      "value": "https://yourdomain.okta.com/oauth2/default/v1/token"
    },
    {
      "key": "@OIDC/USERINFO_URL",
      "value": "https://yourdomain.okta.com/oauth2/default/v1/userinfo"
    },
    {
      "key": "@OIDC/ISSUER",
      "value": "https://yourdomain.okta.com/oauth2/default"
    }
  ]
}
```

## Supported User Information

The OIDC provider extracts the following user information from the UserInfo endpoint:

- **Email**: Primary email address (`email` field)
- **Name**: Full name or constructed from given/family names (`name`, `given_name`, `family_name` fields)
- **Profile Picture**: Avatar/profile image URL (`picture` field)

## Authentication Flow

1. **Authorization**: User is redirected to the OIDC provider's authorization endpoint
2. **Code Exchange**: Authorization code is exchanged for access and refresh tokens
3. **UserInfo Retrieval**: Access token is used to fetch user profile information
4. **Token Refresh**: Refresh tokens are used to obtain new access tokens when needed

## Error Handling

The provider includes comprehensive error handling for:

- Invalid configuration parameters
- Network failures during token exchange
- HTTP error responses from OIDC endpoints
- Malformed JSON responses
- Token expiration and refresh failures

## Security Features

- **PKCE Support**: Includes PKCE (Proof Key for Code Exchange) parameters for enhanced security
- **Token Validation**: Validates HTTP status codes and response formats
- **Secure Headers**: Sets appropriate Accept and Authorization headers
- **Error Sanitization**: Provides clear error messages without exposing sensitive information

## Testing

The package includes comprehensive tests covering:

- Provider configuration and initialization
- Authorization URL generation
- Token exchange and refresh flows
- UserInfo endpoint integration
- Error scenarios and edge cases
- Mock server integration for isolated testing

Run tests with:

```bash
go test ./services/authprovider/oidc/... -v
```

## Compatibility

This OIDC provider is compatible with any OpenID Connect 1.0 compliant identity provider, including:

- Auth0
- Keycloak
- Okta
- Amazon Cognito
- Azure Active Directory (via OIDC endpoints)
- Firebase Auth
- Custom OIDC implementations

The provider follows the OpenID Connect Core 1.0 specification and OAuth 2.0 standards.
