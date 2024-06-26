# ?

This is one out of four parts of [Miku Notes]()

Application parts:
- [Auth service](https://github.com/kuromii5/miku-notes-auth) (this repo)
- [Data service](https://github.com/kutoru/miku-notes-data)
- [Gateway service](https://github.com/kutoru/miku-notes-gateway)
- [Frontend](https://github.com/kinokorain/Miku-notes-frontend)


# How to run

## Preparing

Clone the repo in your project folder:

```bash
git clone https://github.com/kuromii5/miku-notes-auth.git
```

Clone .proto files

```bash
cd miku-notes-auth
git clone https://github.com/kutoru/miku-notes-proto ./proto
```
Make sure you have the [protoc](https://grpc.io/docs/protoc-installation) binary in your path

Run ```task gen``` in root directory to generate code. If you don't have it, you can install task [**here**](https://taskfile.dev/installation/)

## Configuration Setup

To run the project locally, create a directory for your .yaml file, for example:

```bash
mkdir config
```

### Example config

**local.yaml:**

```yaml
env: "local" # dev, prod
postgres:
  user: "postgres"
  password: "admin"
  host: "localhost"
  port: 5432
  dbname: "my_db"
  sslmode: "disable"
token_ttl: 1h
secret: "my_secret"
grpc:
  port: 44044
  timeout: 10h
```

### Migrations

Don't forget to create and run Postgres DB named as in config.

Run migrations using task:
```bash
task migrate
```
You can change migrations table name in Taskfile if you want

## Running app

Adjust **Taksfile.yaml**
```run``` command if needed, then simply run the app:
```
task run
```
You can specify config path if you want by changing this in **Taskfile.yaml**:
```yaml
env:
  CONFIG_PATH: "./config/local.yaml"
```
Or change it directly in cmd:
```bash
task run CONFIG_PATH="my_config/local.yaml"
```

If anything is **ok**, you should receive 2 green log messages.
