# ReserveFlow v1

ReserveFlow v1 is a Go backend project for managing reservable resources such as meeting rooms, desks, courts, or similar bookable assets.

The project includes authentication, admin authorization, resource management, resource-based working hours, and a reservation lifecycle with hold, confirm, and cancel actions.

---

## Tech Stack

* Go
* Gin
* GORM
* PostgreSQL
* JWT
* bcrypt
* godotenv
* Postman
* DBeaver

---

## Main Features

### Auth

* Register user
* Login user
* Generate JWT token
* Get current user profile
* Role-based authorization

### Admin

* Admin protected endpoints
* Admin middleware
* Admin ping endpoint for permission testing

### Resources

* Create resource
* List resources
* Get resource detail
* Update resource
* Delete resource

### Working Hours

* Set working hours for a specific resource
* List working hours for a specific resource
* Resource-based availability rules
* Closed day support

### Reservations

* Hold reservation
* Confirm reservation
* Cancel reservation
* List current user's reservations
* Conflict control
* Working hour control
* Expired hold control

---

## Project Structure

```text
reserveflow-v1/
  api/
    auth.go
    health.go
    resource.go
    reservation.go
    working_hour.go

  commons/
    config.go
    postgres.go

  dao/
    user.go
    resource.go
    reservation.go
    working_hour.go

  middleware/
    auth.go
    role.go

  models/
    user.go
    resource.go
    reservation.go
    working_hour.go

  .env
  go.mod
  main.go
  README.md
```

---

## Environment Variables

Create a `.env` file in the project root:

```env
APP_NAME=reserveflow-v1
APP_ENV=development
APP_PORT=8083

POSTGRES_HOST=localhost
POSTGRES_PORT=5433
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=reserveflow_v1
POSTGRES_SSLMODE=disable

JWT_SECRET=change-this-secret
JWT_ACCESS_TTL_MINUTES=15
```

---

## Database Setup

Create the database in PostgreSQL:

```sql
CREATE DATABASE reserveflow_v1;
```

The project uses GORM AutoMigrate for creating tables.

Current tables:

```text
users
resources
working_hours
reservations
```

---

## How to Run

Install dependencies:

```bash
go mod tidy
```

Run the project:

```bash
go run main.go
```

Expected output:

```text
Successfully connected to database
Server running on port: 8083
```

Base URL:

```text
http://localhost:8083
```

---

## Authorization

Protected endpoints require a JWT token in the `Authorization` header.

This project currently uses raw token format:

```text
Authorization: YOUR_TOKEN
```

Bearer format is not used yet.

---

## API Endpoints

### Health

| Method | Endpoint  | Description                 |
| ------ | --------- | --------------------------- |
| GET    | `/health` | Check if the API is running |

---

### Auth

| Method | Endpoint         | Description              |
| ------ | ---------------- | ------------------------ |
| POST   | `/auth/register` | Register a new user      |
| POST   | `/auth/login`    | Login and get token      |
| GET    | `/auth/me`       | Get current user profile |

#### Register Request

```json
{
  "full_name": "Girayhan Durmus",
  "email": "girayhan@reserveflow.com",
  "password": "123456"
}
```

#### Login Request

```json
{
  "email": "girayhan@reserveflow.com",
  "password": "123456"
}
```

---

### Admin

| Method | Endpoint      | Description       |
| ------ | ------------- | ----------------- |
| GET    | `/admin/ping` | Test admin access |

This endpoint requires an admin token.

---

### Resources

| Method | Endpoint               | Description         | Auth   |
| ------ | ---------------------- | ------------------- | ------ |
| GET    | `/resources`           | List all resources  | Public |
| GET    | `/resources/:id`       | Get resource detail | Public |
| POST   | `/admin/resources`     | Create resource     | Admin  |
| PATCH  | `/admin/resources/:id` | Update resource     | Admin  |
| DELETE | `/admin/resources/:id` | Delete resource     | Admin  |

#### Create Resource Request

```json
{
  "name": "Meeting Room A",
  "description": "8 kişilik toplantı odası",
  "capacity": 8
}
```

#### Update Resource Request

```json
{
  "name": "Meeting Room A Updated",
  "description": "10 kişilik güncellenmiş toplantı odası",
  "capacity": 10,
  "is_active": true
}
```

---

### Working Hours

| Method | Endpoint                             | Description                       | Auth  |
| ------ | ------------------------------------ | --------------------------------- | ----- |
| POST   | `/admin/resources/:id/working-hours` | Set working hours for a resource  | Admin |
| GET    | `/admin/resources/:id/working-hours` | List working hours for a resource | Admin |

#### Set Working Hours Request

```json
{
  "working_hours": [
    {
      "day_of_week": "monday",
      "open_time": "09:00",
      "close_time": "18:00",
      "is_closed": false
    },
    {
      "day_of_week": "tuesday",
      "open_time": "09:00",
      "close_time": "18:00",
      "is_closed": false
    },
    {
      "day_of_week": "wednesday",
      "open_time": "09:00",
      "close_time": "18:00",
      "is_closed": false
    },
    {
      "day_of_week": "thursday",
      "open_time": "09:00",
      "close_time": "18:00",
      "is_closed": false
    },
    {
      "day_of_week": "friday",
      "open_time": "09:00",
      "close_time": "18:00",
      "is_closed": false
    },
    {
      "day_of_week": "saturday",
      "open_time": "10:00",
      "close_time": "16:00",
      "is_closed": false
    },
    {
      "day_of_week": "sunday",
      "open_time": "",
      "close_time": "",
      "is_closed": true
    }
  ]
}
```

---

### Reservations

| Method | Endpoint                    | Description                      | Auth |
| ------ | --------------------------- | -------------------------------- | ---- |
| POST   | `/reservations/hold`        | Hold a reservation slot          | User |
| POST   | `/reservations/:id/confirm` | Confirm a held reservation       | User |
| POST   | `/reservations/:id/cancel`  | Cancel a reservation             | User |
| GET    | `/reservations/my`          | List current user's reservations | User |

#### Hold Reservation Request

```json
{
  "resource_id": 2,
  "start_time": "2026-06-17T10:00:00+03:00",
  "end_time": "2026-06-17T11:00:00+03:00"
}
```

#### Confirm Reservation

```text
POST /reservations/2/confirm
```

No request body is required.

#### Cancel Reservation

```text
POST /reservations/2/cancel
```

No request body is required.

---

## Reservation Status Flow

Reservation statuses:

```text
held
confirmed
cancelled
expired
```

Main flow:

```text
held -> confirmed
held -> cancelled
confirmed -> cancelled
held -> expired
```

---

## Business Rules

### Resource Rules

* Reservation can only be created for an existing resource.
* Resource must be active.
* Deleted resources cannot be used.

### Working Hour Rules

* Reservation time must be inside the resource working hours.
* If the selected day is closed, reservation is rejected.
* If working hour record does not exist for the selected day, reservation is rejected.
* Time format for working hours is `HH:mm`.

### Reservation Rules

* `end_time` must be after `start_time`.
* `start_time` and `end_time` must be RFC3339 format.
* A reservation starts as `held`.
* A held reservation expires after 10 minutes.
* A held reservation can be confirmed before expiration.
* Expired held reservations cannot be confirmed.
* Cancelled reservations cannot be cancelled again.
* Confirmed reservations block the time slot.
* Non-expired held reservations also block the time slot.
* Cancelled and expired reservations do not block the time slot.

---

## Postman Collection Structure

Recommended Postman folder structure:

```text
ReserveFlow-v1
  Auth
    Health Check
    Register User
    Login User
    Get Me

  Admin
    Admin Ping

  Resource
    Create Resource
    Get All Resources
    Get Resource By ID
    Update Resource
    Delete Resource

  Working Hour
    Set Working Hours
    Get Working Hours

  Reservation
    Hold Reservation
    Confirm Reservation
    Cancel Reservation
    Get My Reservations
```

---

## Suggested Test Flow

1. Health Check
2. Register User
3. Login User
4. Get Me
5. Change user role to `admin` in database
6. Login again and get admin token
7. Admin Ping
8. Create Resource
9. Get All Resources
10. Set Working Hours
11. Get Working Hours
12. Hold Reservation
13. Try same Hold Reservation again and expect conflict
14. Confirm Reservation
15. Try same Confirm Reservation again and expect error
16. Cancel Reservation
17. Try same Cancel Reservation again and expect error
18. Get My Reservations

---

## Example Error Response

```json
{
  "success": false,
  "error": {
    "code": "RESERVATION_CONFLICT",
    "message": "resource is already reserved or held for this time range"
  }
}
```

---

## Future Improvements

Possible next improvements:

* Add `.env.example`
* Add Docker Compose
* Add unit tests
* Add integration tests
* Add structured logging
* Add request validation
* Add pagination
* Add blackout dates
* Add background worker for expiring held reservations
* Add service layer
* Add Swagger/OpenAPI documentation
* Add Bearer token support
* Add database migrations instead of AutoMigrate

---

## Current Status

ReserveFlow v1 is a learning-focused backend project.

The core reservation lifecycle is completed:

```text
resource -> working hours -> hold -> confirm -> cancel
```

The project is ready for documentation, testing, Dockerization, and future production-oriented refactoring.
