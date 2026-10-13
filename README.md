# Online Course Management System

A REST API for managing courses, students, and course enrollments. The application uses Go, Gin, GORM, and PostgreSQL.

## Features

- Create, list, retrieve, update, and delete courses and students.
- Search and paginate course and student lists.
- Enroll a student in a course.
- Enforce a maximum of 3 courses per student and 10 students per course.
- Reject duplicate enrollments.
- List the five courses with the most enrolled students.
- List students not enrolled in any course.

## Requirements

- Go 1.27 or newer
- PostgreSQL

## Setup

1. Create a database:

   ```sql
   CREATE DATABASE course_management;
   ```

2. Copy the example environment file:

   ```powershell
   Copy-Item .env.example .env
   ```

3. Update `POSTGRES_DSN` in `.env` with your PostgreSQL credentials:

   ```env
   POSTGRES_DSN="host=localhost user=postgres password=your_password dbname=course_management port=5432 sslmode=disable"
   ```

4. Download dependencies and start the API:

   ```powershell
   go mod download
   go run ./cmd/course-management
   ```

The server listens at `http://localhost:8080`. On startup, GORM automatically creates or updates the `courses`, `students`, and `course_students` tables.

### Background summary job

The application starts a background job that logs the total courses, total students, and average course duration. It currently runs every 10 seconds to make local verification easy. For the assignment's required two-minute interval, change the ticker duration in `internal/scheduler/summary_scheduler.go` to:

```go
time.NewTicker(2 * time.Minute)
```

Press `Ctrl+C` to stop both the scheduler and HTTP server gracefully.

## API

All endpoints are prefixed with `/api/v1`.

### Courses

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/courses` | List courses. Supports `limit`, `offset`, and `search`. |
| `GET` | `/courses/:id` | Get one course, including its students. |
| `POST` | `/courses` | Create a course. |
| `PUT` | `/courses/:id` | Update a course. Send only fields to change. |
| `DELETE` | `/courses/:id` | Delete a course. |
| `POST` | `/courses/:id/enroll` | Enroll a student in a course. |
| `GET` | `/courses/top-enrolled` | Get up to five courses ordered by enrollment count. |

Create a course:

```powershell
curl.exe -X POST http://localhost:8080/api/v1/courses `
  -H "Content-Type: application/json" `
  -d '{"title":"Go Fundamentals","description":"An introduction to Go.","durationInHours":12,"startDate":"2026-11-01T00:00:00Z","instructor":"Ada Lovelace"}'
```

List courses containing `go` in the title:

```text
GET /api/v1/courses?search=go&limit=20&offset=0
```

Update a course:

```powershell
curl.exe -X PUT http://localhost:8080/api/v1/courses/1 `
  -H "Content-Type: application/json" `
  -d '{"durationInHours":16}'
```

### Students

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/students` | List students. Supports `limit`, `offset`, and `search`. |
| `GET` | `/students/:id` | Get one student, including enrolled courses. |
| `POST` | `/students` | Create a student. |
| `PUT` | `/students/:id` | Update a student's name or email. |
| `DELETE` | `/students/:id` | Delete a student. |
| `GET` | `/students/unenrolled` | List students with no enrolled courses. |

Create a student:

```powershell
curl.exe -X POST http://localhost:8080/api/v1/students `
  -H "Content-Type: application/json" `
  -d '{"name":"Grace Hopper","email":"grace@example.com"}'
```

List students whose name contains `grace`:

```text
GET /api/v1/students?search=grace&limit=20&offset=0
```

### Enrollment

Enroll student `1` in course `1`:

```powershell
curl.exe -X POST http://localhost:8080/api/v1/courses/1/enroll `
  -H "Content-Type: application/json" `
  -d '{"studentID":1}'
```

The endpoint returns:

- `200 OK` when enrollment succeeds.
- `400 Bad Request` when the course already has 10 students or the student already has 3 courses.
- `409 Conflict` when the student is already enrolled in that course.
- `404 Not Found` when the course or student does not exist.

Enrollments cannot be changed through `PUT /students/:id`; use the enrollment endpoint so the business rules are always enforced.

## Pagination and search

Both list endpoints accept:

| Parameter | Meaning | Default / limit |
| --- | --- | --- |
| `limit` | Maximum records returned | Defaults to 20; maximum 100 |
| `offset` | Number of records to skip | Defaults to 0 |
| `search` | Case-insensitive name or title substring | Optional |

Examples:

```text
GET /api/v1/courses?limit=10&offset=20
GET /api/v1/students?search=hopper
```

## Validation and errors

- A course requires `title`, `durationInHours`, and `instructor`.
- A student requires `name` and a valid `email`.
- Course titles must be non-empty and at most 200 characters.
- Invalid JSON or query parameters return `400 Bad Request`.
- Duplicate course titles and duplicate enrollments return `409 Conflict`.

Error responses use this format:

```json
{
  "error": "validation error: a student cannot enroll in more than 3 courses"
}
```

## Project layout

```text
cmd/course-management/    Application entry point
internal/config/          Database connection setup
internal/handler/         HTTP handlers and routes
internal/model/           Entities and request/response models
internal/repository/      GORM data access
internal/service/         Business logic
internal/mapper/          API response mapping
```

## Verify the build

```powershell
go test ./...
```
