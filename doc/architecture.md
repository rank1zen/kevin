# Architecture

## Project Structure

Only significant patterns are documented

```
kevin/
├── cmd/kevin/             # Main web server entry point
├── internal/              # Private application code
│   ├── feature/<feature>/ # Per feature packages
│   ├── web/               # Web related
│   │   ├── templui/       # Reusable components (templui)
│   │   ├── component/     # Reusable components
│   │   ├── block/         # Reusable UI blocks
│   │   ├── partial/       # Loading UI blocks
│   │   └── page/          # Full pages
│   ├── db/                # Repository implementations
│   ├── postgres/          # Postgres setup & migrations
│   ├── validation/        # Input validators
│   ├── config/            # Configuration management
│   ├── ctxkeys/           # Context key definitions
│   ├── log/               # Structured logging
│   ├── app/               # App initialization
│   ├── route/             # Route definitions
│   └── utils/             # Shared utilities
└── assets/                # Static files (CSS, JS, images)
```
