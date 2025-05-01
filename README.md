# 🏱 Lucky Draw - Backend

A backend system for lucky draw activities and prize distribution, built with **Go + Gin**, using **Bun ORM** and **PostgreSQL**. The system supports player import, prize management, conditional drawing, and winner tracking through RESTful APIs.

---

## 📌 Features

- ✅ Manage event rooms
- ✅ Import players via `.csv`
- ✅ Manage prizes with Cloudinary image upload
- ✅ Define draw conditions based on player positions
- ✅ Draw winners based on conditions
- ✅ Store and query winner history
- ✅ List all related data by room (ListAll API)
- ✅ Soft delete supported on all models
- ✅ String-based UUIDs for easier handling

---

## 🏗️ Tech Stack

| Layer       | Technology                    |
|-------------|-------------------------------|
| Backend     | Go (Golang)                   |
| Framework   | Gin                           |
| ORM         | Bun, Npm                      |
| Database    | PostgreSQL                    |
| File Upload | Cloudinary                    |
| OS          | Linux (Fedora, Rocky)         |
| API Format  | RESTful JSON APIs             |

---

## 📂 Project Structure Example

```
Lucky-Draw_Back-end/
|
├── app/
│   ├── controller/       # API handlers
│   ├── service/          # Business logic
│   ├── model/            # ORM definitions
│   ├── request/          # Input schemas
│   ├── response/         # Output schemas
│
├── config/               # Environment and config
├── internal/logger/      # Logging utility
├── utils/                # Helper functions (CSV, etc.)
├── main.go               # Entrypoint
└── go.mod / go.sum       # Go modules
```

---

## 📏 Models Overview

### Room
- `id`, `name`, `created_at`, `updated_at`, `deleted_at`

### Player
- `id`, `prefix`, `first_name`, `last_name`, `member_id`, `position`, `room_id`, `is_active`, `status`, `created_at`, `updated_at`, `deleted_at`

### Prize
- `id`, `name`, `image_url`, `quantity`, `room_id` ,`created_at`, `updated_at`, `deleted_at`

### DrawCondition
- `id`, `room_id`, `prize_id`, `filter_status`, `filter_is_active`, `filter_position`, `quantity`, `created_at`, `updated_at`, `deleted_at`

### Winner
- `id`, `room_id`, `player_id`, `prize_id`, `draw_condition_id` ,`created_at`, `updated_at`, `deleted_at`

---

## 📊 API Endpoints

### Room
| Method | Endpoint                    | Description               |
|--------|-----------------------------|---------------------------|
| GET    | `/api/v1/rooms`             | List rooms                |
| POST   | `/api/v1/rooms`             | Create a new room         |
| PUT    | `/api/v1/rooms/:id`         | Update room               |
| DELETE | `/api/v1/rooms/:id`         | Soft delete room          |
| GET    | `/api/v1/rooms/:id/all`     | Fetch all room-related data |

### Player
| Method | Endpoint                              | Description              |
|--------|----------------------------------------|--------------------------|
| POST   | `/api/v1/players/import/:roomId`       | Import from CSV          |
| GET    | `/api/v1/players/room/:roomId`         | List players by room     |

### Prize
| Method | Endpoint                              | Description              |
|--------|----------------------------------------|--------------------------|
| GET    | `/api/v1/prizes/room/:roomId`         | List prizes by room      |
| POST   | `/api/v1/prizes`                      | Create prize             |
| PUT    | `/api/v1/prizes/:id`                  | Update prize             |
| DELETE | `/api/v1/prizes/:id`                  | Delete prize             |

### DrawCondition
| Method | Endpoint                                                      | Description                  |
|--------|----------------------------------------------------------------|------------------------------|
| POST   | `/api/v1/draw-conditions`                                     | Create draw condition        |
| GET    | `/api/v1/draw-conditions/prize/:prizeId`                      | List conditions for a prize  |
| GET    | `/api/v1/draw-conditions/GetDrawConditionPreview/:conditionId`| Preview matched players      |

### Winner
| Method | Endpoint                            | Description              |
|--------|--------------------------------------|--------------------------|
| POST   | `/api/v1/winners/draw`              | Draw a winner            |
| GET    | `/api/v1/winners/room/:roomId`      | List winners in room     |

---

## 📄 CSV Import Format

- CSV must be without a header row
- Column order: `prefix`, `first_name`, `last_name`, `member_id`, `position`

**Example:**
```
Mr.,John,Doe,CS001,Student
Ms.,Anna,Smith,CS002,Staff
```

---

## ☁️ Cloudinary Setup

Set the following environment variables for image uploads:

```env
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
```

---

## 📅 .env Example

```env
DATABASE_URL=postgres://user:password@localhost:5432/lucky_draw_db?sslmode=disable
PORT=8080

CLOUDINARY_CLOUD_NAME=xxx
CLOUDINARY_API_KEY=xxx
CLOUDINARY_API_SECRET=xxx
```

---

## ▶️ Getting Started

### 1. Clone the repository
```bash
git clone https://github.com/pathipat11/Lucky-Draw_Back-end.git
cd Lucky-Draw_Back-end
```

### 2. Create a `.env` file
Fill in your database and Cloudinary credentials

### 3. Run the server
```bash
go run main.go
```

---

## 📩 Contact

Developed by **Palmy (Mata)**  
Email: pathipat.mattra@gmail.com  
GitHub: [github.com/pathipat11](https://github.com/pathipat11)

---
