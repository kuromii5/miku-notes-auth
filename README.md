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
  refresh_ttl: 720h
  redis_addr: "127.0.0.1:6379"
  secret: "my_secret"
grpc:
  port: 44044 # port for your gRPC server
  timeout: 10h # max duration of time that the request can take
```

### Migrations

Don't forget to create and run Postgres DB named as in config.

Run migrations with `task migrate` or using cmd, for example:

```bash
go run cmd/migrations/main.go --migrations-table="migrations" --config="config/local.yaml"
```

## Running app

Simply run the app:

```bash
task run CONFIG_PATH="./my_config/local.yaml"
```

Or if using cmd:

```bash
go run cmd/sso/main.go --config="./my_config/local.yaml"
```

If anything is **ok**, you should receive 2 green log messages.
