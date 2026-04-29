# Architecture

## Project Structure

*Only significant patterns are documented*

```
kevin/
├── main.go                    # Entry point
├── cmd/                       # Cobra CLI application
├── internal/                  # Private application code
│   ├── feature/               # Per-feature packages
│   │   └── <feature>/         # Feature-specific package
│   │       ├── ui/            # Server-rendered UI
│   │       │   ├── handler/   # Handlers for rendering UI
│   │       │   ├── component/ # Reusable components
│   │       │   ├── block/     # Reusable blocks
│   │       │   ├── skeleton/  # Loading skeletons
│   │       │   ├── layout/    # Page layouts
│   │       │   └── page/      # Full pages
│   │       ├── api/           # Server API
│   │       │   └── handler/   # Handlers for serving JSON
│   │       ├── service/       # Business logic
│   │       ├── model/         # Data models
│   │       ├── db/            # Store implementations
│   │       └── route/         # Route definitions
│   ├── postgres/              # Postgres settings and migrator
│   ├── config/                # Configuration management
│   ├── log/                   # Structured logging
│   ├── app/                   # App initialization
│   └── route/                 # Route definitions
├── migrations/                # Embedded migrations
└── assets/                    # Static files (CSS, JS, images)
```
