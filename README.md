## Stack Seeker

## Install these tools:
[sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
[sql-migrate](https://github.com/rubenv/sql-migrate)

### Install tldr;
```bash
$ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
$ go install github.com/rubenv/sql-migrate/...@latest
```

## Create a new migration run:
```bash
$ sql-migrate new {name}
```

## Generate the SQLC code:
```bash
$ go generate ./...
```
