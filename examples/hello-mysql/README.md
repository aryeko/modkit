# hello-mysql

Example consuming app for modkit using MySQL, sqlc, and migrations.

## Run

```bash
make run
```

Then hit:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/users/1
```

## Test

```bash
make test
```
