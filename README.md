## Go / Gin / PostgreSQL Backend Boilerplate

### Core

 - Go + Gin

 - PostgreSQL

 - JWT authentication

 - Configuration management

 - REST API

 - Environment-based configuration

### Authentication & Authorization

 - Login / Logout

 - JWT access & refresh tokens

 - Token refresh / session management

 - Password hashing

 - Role-based access control

 - Roles: `Admin`, `Customer`

### Customer

 - Register

    -  Username

    -  Email

    -  Alternative email (optional)

    -  Password

 - Login with username/email

 - Logout

 - View/update profile

 - Change password

 - Refresh session

### Admin

 - View/update own profile

 - Change own password

 - Create customer

 - Update customer

 - Search customers

 - Get customer by ID

 - Change customer status

 - Change customer password

### Database

 - PostgreSQL connection

 - Connection management

 - Indexes

 - Repository layer

### Security

 - Password hashing

 - Access/refresh token rotation

 - Input validation

 - Rate limiting

 - CORS

 - Authentication & authorization middleware

### Boilerplate Essentials

 - Standard API response/error format

 - Request validation

 - Pagination & filtering

 - Logging

 - Health check

 - Graceful shutdown

`.env.example`


### Login Identifier Settings

Administrators can control which identifiers users are allowed to use when signing in.

Available options:

- **Primary Email** — Allow users to log in using their primary email address.
- **Username** — Allow users to log in using their username.
- **Phone Number** — Allow users to log in using their registered 10-digit phone number.
- **Alternate Email** — Allow users to log in using their alternate email address.

At least **one login identifier must remain enabled** at all times. Disabling an identifier prevents users from signing in with that identifier but does not affect the user's stored account information.

These settings apply globally to user login.

### .env

```
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<db_password>
DB_NAME=<db_name>
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Kolkata

# Server
PORT=<serverport>
GIN_MODE=debug
# or release

# Authentication
JWT_SECRET=<secret_long_string>
JWT_EXPIRY_HOURS=24

AUTH_ACCESS_COOKIE=access_token
AUTH_REFRESH_COOKIE=refresh_token

REFRESH_TOKEN_EXPIRY_DAYS=30
REFRESH_TOKEN_REVOKED_RETENTION_DAYS=7
REFRESH_TOKEN_EXPIRED_RETENTION_DAYS=7

# Cookies
COOKIE_SECURE=false
COOKIE_SAME_SITE=lax

# CORS
ALLOWED_ORIGINS=http://localhost:<frontend-serverport>

# Two-factor authentication
TOTP_ISSUER=Basic App
TOTP_ENCRYPTION_KEY=<32-byte-secret>

# Backup codes
BACKUP_CODE_COUNT=10
BACKUP_CODE_LOW_THRESHOLD=4

#Seed
SEED_ADMIN_FIRST_NAME=System
SEED_ADMIN_LAST_NAME=Administrator
SEED_ADMIN_USERNAME=admin
SEED_ADMIN_EMAIL=admin@template.local
SEED_ADMIN_PHONE=0000000000
SEED_ADMIN_PASSWORD=change-this-password

```


## Migrations

```
go run ./cmd/migrate up
```

```
go run ./cmd/migrate version

```

## Seed Account

```
go run ./cmd/seed
```

## Run air

If air is setup correctly. Use the below command to run the application.

```
ENV=development air

```

## Getting Started

Run the development server:

```
npm run dev

```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.


### Other migration commands


```
go run ./cmd/migrate down
```

```
go run ./cmd/migrate force <version>
```
```
go run ./cmd/migrate force 1
```

