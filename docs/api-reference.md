# Health Endpoints

| Method | Endpoint     | Description         |
| ------ | ------------ | ------------------- |
| GET    | `/health`    | Simple health check |
| GET    | `/readiness` | Readiness probe     |
| GET    | `/liveness`  | Liveness probe      |

## Auth Routes (`/api/v1`)

| Method | Endpoint           | Description         |
| ------ | ------------------ | ------------------- |
| POST   | `/register`        | Register new user   |
| POST   | `/login`           | Login user          |
| POST   | `/password/forgot` | Request reset token |
| POST   | `/password/reset`  | Reset user password |

### Profile Routes (JWT required)

| Method | Endpoint                   | Description         |
| ------ | -------------------------- | ------------------- |
| GET    | `/profile`                 | Get own profile     |
| PUT    | `/profile`                 | Update own profile  |
| DELETE | `/profile`                 | Delete own account  |
| PUT    | `/profile/change-password` | Change own password |
| POST   | `/profile/upload-avatar`   | Upload avatar image |

### Admin Routes (`/api/v1/admin`)

(Require `manage_users`, `manage_teams`, etc.)

- **Users**: CRUD for all users, permission update/reset
- **Teams**: CRUD, upload logo
- **Seasons**: CRUD, upload logo
- **Matches**: CRUD, upload video & scout file
- **Audit Logs**: View admin actions
- **Waitlist**: Approve/Reject

---

### Match Endpoint with JSON Parser Integration

| Method | Endpoint             | Description                              |
| ------ | -------------------- | ---------------------------------------- |
| GET    | `/admin/matches/:id` | Returns match details + parsed JSON data |
