seo-backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── db.go
│   │   └── migrations/
│   │       ├── 001_create_users_table.up.sql
│   │       ├── 001_create_users_table.down.sql
│   │       ├── 002_create_teams_table.up.sql
│   │       ├── 002_create_teams_table.down.sql
│   │       ├── 003_create_products_table.up.sql
│   │       ├── 003_create_products_table.down.sql
│   │       ├── 004_create_adapter_configs_table.up.sql
│   │       ├── 004_create_adapter_configs_table.down.sql
│   │       ├── 005_create_drafts_table.up.sql
│   │       ├── 005_create_drafts_table.down.sql
│   │       ├── 006_create_histories_table.up.sql
│   │       ├── 006_create_histories_table.down.sql
│   │       ├── 007_create_api_keys_table.up.sql
│   │       ├── 007_create_api_keys_table.down.sql
│   │       └── 008_create_settings_table.up.sql
│   │       └── 008_create_settings_table.down.sql
│   ├── models/
│   │   ├── user.go
│   │   ├── team.go
│   │   ├── product.go
│   │   ├── draft.go
│   │   └── history.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── draft.go
│   │   ├── history.go
│   │   └── generate.go
│   ├── middleware/
│   │   └── auth.go
│   └── repository/
│       ├── user_repo.go
│       ├── product_repo.go
│       └── draft_repo.go
├── go.mod
├── go.sum
├── .env
└── Makefile