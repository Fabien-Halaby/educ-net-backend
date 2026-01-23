```sh
educ-net-backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   └── postgres.go
│   ├── domain/
│   │   ├── error.go
│   │   ├── school.go
│   │   └── user.go
│   ├── repository/
│   │   ├── school_repository.go
│   │   └── user_repository.go
│   ├── usecase/
│   │   └── school_usecase.go
│   ├── handler/
│   │   ├── dto/
│   │   │   └── school_dto.go
│   │   └── school_handler.go
│   ├── middleware/
│   │   ├── middleware.go
|   |   ├── cors.go
|   |   └── logging.go
│   └── utils/
│       ├── response.go
│       └── slug.go
├── migrations/
│   └── 001_init.sql
├── .env
├── go.mod
├── README.md
└── STRUCTURE.md
```