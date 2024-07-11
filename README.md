# ?

This is one out of four parts of [Miku Notes]()

Application parts:

- [Auth service](https://github.com/kuromii5/miku-notes-auth) (this repo)
- [Data service](https://github.com/kutoru/miku-notes-data)
- [Gateway service](https://github.com/kutoru/miku-notes-gateway)
- [Frontend](https://github.com/kinokorain/Miku-notes-frontend)

# How to run

## Preparing

Clone the repo and corresponding submodule for generating server code into your folder

Make sure you have the [protoc](https://grpc.io/docs/protoc-installation) binary in your path

Run `task gen` in root directory to generate code. If you don't have it, you can install task [**here**](https://taskfile.dev/installation/) **or** use cmd described in **Taskfile.yaml**

## Configuration Setup

### Example config

**local.yaml:**

```yaml
env: "local" # environment can be "dev", "prod" or "local"
postgres: # postgres settings
  user: "postgres"
  password: "admin"
  host: "localhost"
  port: 5432
  dbname: "my_db"
  sslmode: "disable"
tokens:
  access_ttl: 15m
  refresh_ttl: 720h # 30 days
  redis_addr: "127.0.0.1:6379" # address of redis db which stores refresh tokens
  secret: "my_secret" # for JWT access tokens
grpc:
  port: 44044 # port for your gRPC server
  connection_token: "private_connection_token" # auth token for secure connection between services
```

### OR

**.env:**

```dotenv
# ENVIRONMENT
ENV=local # dev, prod

# POSTGRES SETTINGS
POSTGRES_USER=postgres
POSTGRES_PASSWORD=admin
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DBNAME=miku_notes
POSTGRES_SSLMODE=disable

# TOKEN MANAGEMENT SETTINGS
TOKENS_ACCESS_TTL=15m
TOKENS_REFRESH_TTL=720h
TOKENS_REDIS_ADDR=127.0.0.1:6379
TOKENS_SECRET=my_secret

# GPRC SETTINGS
GRPC_CONNECTION_TOKEN=private_connection_token
GRPC_PORT=44044
```

### Migrations

Don't forget to create and run Postgres DB named as in config.

Run migrations with `task migrate` or using cmd (if using .env - don't specify cfg path), for example:

```bash
go run cmd/migrations/main.go --migrations-table="migrations" --config="config/local.yaml"
```

## Running app

Simply run the app:

```bash
task run CONFIG_PATH="./my_config/local.yaml"
```

Or if using cmd, if you are using .yaml config, you also need to provide config path:

```bash
go run cmd/sso/main.go --config="./my_config/local.yaml"
```

If anything is **ok**, you should receive 2 green log messages.
