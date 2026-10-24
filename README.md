# เจ้าพ่อศรีนครเตา — Backend (Go API)

JSON API, admin auth, RAG chat and uploads for the faith-tourism platform. The frontend
([FE-Nakrontao](https://github.com/puyolnw/FE-Nakrontao)) proxies `/api` and `/uploads` here, so the
browser only ever talks to one origin.

## Run locally

```bash
cp .env.example .env   # fill DATABASE_URL, IP_HASH_SALT, BOOTSTRAP_ADMIN_*
go run ./cmd/api
```

The listen address comes from `ADDR`, else `PORT`, else `:8080`.

## Deploy on Railway

1. Add a Postgres service **with pgvector** (the schema runs `CREATE EXTENSION vector` and `pg_trgm`).
2. Add this repo as a service; it builds from `Dockerfile`. Railway injects `PORT`.
3. Set variables from `.env.example` (at least `DATABASE_URL`, `IP_HASH_SALT`, `BOOTSTRAP_ADMIN_EMAIL`,
   `BOOTSTRAP_ADMIN_PASSWORD`, `TRUSTED_PROXY=fd00::/8`).
4. Mount a volume at `/data/uploads`.
5. Keep this service private (no public domain); the frontend reaches it over Railway's private network,
   e.g. `http://be-nakrontao.railway.internal:8080`.

Migrations run on start. `SEED_FILE` loads content once into an empty database.
