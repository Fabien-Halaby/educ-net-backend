```sh
educ-net-backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── models/                  # Data models
│   │   ├── school.go
│   │   ├── user.go
│   │   ├── channel.go
│   │   ├── post.go
│   │   └── notification.go
│   ├── handlers/                # HTTP handlers
│   │   ├── auth.go
│   │   ├── schools.go
│   │   ├── users.go
│   │   ├── channels.go
│   │   ├── posts.go
│   │   ├── notifications.go
│   │   └── upload.go
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── database/                # Database layer
│   │   └── postgres.go
│   └── websocket/               # WebSocket hub
│       └── hub.go
├── migrations/                  # SQL migrations
│   └── 001_init.sql
├── uploads/                     # Uploaded files (gitignored)
├── .env.example                 # Example environment config
├── .gitignore
├── go.mod
├── go.sum
├── README.md
└── STRUCTURE.md
```