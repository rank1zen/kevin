# Kevin

A lightweight League of Legends stats client.

<img width="1330" height="1080" alt="combined" src="https://github.com/user-attachments/assets/a45ba9f4-72cc-4446-a1ab-895fd733891c" />

## Quick Start

Deploy Kevin locally using Docker for early testing and development:

### Prerequisites

*   **[Docker](https://docs.docker.com/get-docker/)**: Installed and running on your machine.
*   **[Go (1.20+)](https://go.dev/doc/install)**: Installed locally if you want to run database migrations directly from your host.

### 1. Obtain a Riot Games API Key

Kevin requires a Riot Games API key to fetch data.
1.  Visit the [Riot Games Developer Portal](https://developer.riotgames.com/).
2.  Sign in and generate a Development API Key.
3.  Keep this key safe; you'll use it as an environment variable.

### 2. Build the Docker Image

Navigate to the root of the `kevin` project and build the Docker image:

```bash
docker build -t kevin-frontend:latest .
```

### 3. Start a local PostgreSQL Database

Kevin uses PostgreSQL for data storage. Start a local PostgreSQL instance using Docker:

```bash
# Create data directory (optional, Docker will create if not exists)
# For quick local start, a named volume is simpler.
docker volume create kevin-db-data

docker run -d \
  --name kevin-postgres \
  -e POSTGRES_USER=kevin \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=kevin \
  -p 5432:5432 \
  --mount source=kevin-db-data,target=/var/lib/postgresql/data \
  postgres:16
```

Wait a few moments for the database to initialize and start up. You can check its status with `docker logs kevin-postgres`.

### 4. Run Database Migrations

Apply the database migrations to your local PostgreSQL instance. You'll need the `tern` tool, which is used in this project. You can run it locally:

```bash
# Ensure you are in the root directory of the kevin project
export KEVIN_POSTGRES_CONNECTION="postgresql://kevin:password@localhost:5432/kevin?sslmode=disable"
go run github.com/jackc/tern/v2 migrate -c tern.conf
```
*Note: The `KEVIN_POSTGRES_CONNECTION` environment variable needs to be set for the `tern` command to connect to your PostgreSQL container. Here, `localhost` is used as `tern` runs on your host machine.*

### 5. Run the Kevin Application Container

Now, run the Kevin frontend application, passing the required environment variables:

```bash
docker run -d \
  --name kevin-app \
  -p 4001:4001 \
  -e KEVIN_RIOT_API_KEY="YOUR_RIOT_API_KEY_HERE" \
  -e KEVIN_POSTGRES_CONNECTION="postgresql://kevin:password@kevin-postgres:5432/kevin?sslmode=disable" \
  -e KEVIN_ENV="development" \
  --network host \
  kevin-frontend:latest
```
*Replace `YOUR_RIOT_API_KEY_HERE` with your actual Riot Games API key.*
*Note: `--network host` is used to allow the `kevin-app` container to connect to the `kevin-postgres` container using `localhost` on some Docker setups (e.g., Docker Desktop on macOS/Windows). If running on Linux, you might prefer to use a Docker network and connect directly via the container name `kevin-postgres`.*

### 6. Access the Application

Once the `kevin-app` container is running (check `docker logs kevin-app`), access the application in your web browser at:

`http://localhost:4001`
