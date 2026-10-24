CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Text fields that need translation are stored as JSONB: {"th": "...", "zh": "...", "en": "..."}

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    name          TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL CHECK (role IN ('admin', 'editor')),
    active        BOOLEAN NOT NULL DEFAULT TRUE,
    failed_logins INT NOT NULL DEFAULT 0,
    locked_until  TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    token_hash BYTEA PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX sessions_user_idx ON sessions(user_id);

CREATE TABLE audit_log (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users(id) ON DELETE SET NULL,
    action     TEXT NOT NULL,
    entity     TEXT NOT NULL,
    entity_id  TEXT NOT NULL,
    detail     JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE site_settings (
    key   TEXT PRIMARY KEY,
    value JSONB NOT NULL
);

-- Sheet "Knowledge"
CREATE TABLE knowledge (
    id         BIGSERIAL PRIMARY KEY,
    doc_code   TEXT NOT NULL UNIQUE,
    category   TEXT NOT NULL,
    kind       TEXT NOT NULL CHECK (kind IN ('history', 'belief', 'ritual', 'practice', 'community')),
    title      JSONB NOT NULL,
    body       JSONB NOT NULL,
    source     TEXT NOT NULL DEFAULT '',
    informant  TEXT NOT NULL DEFAULT '',
    reviewer   TEXT NOT NULL DEFAULT '',
    status     TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'review', 'published')),
    cover_url  TEXT NOT NULL DEFAULT '',
    video_url  TEXT NOT NULL DEFAULT '',
    sort       INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Sheet "Places"
CREATE TABLE places (
    id            BIGSERIAL PRIMARY KEY,
    code          TEXT NOT NULL UNIQUE,
    category      TEXT NOT NULL CHECK (category IN ('sacred', 'culture', 'history', 'community', 'food', 'stay')),
    name          JSONB NOT NULL,
    description   JSONB NOT NULL,
    lat           DOUBLE PRECISION NOT NULL,
    lon           DOUBLE PRECISION NOT NULL,
    open_hours    TEXT NOT NULL DEFAULT '',
    visit_minutes INT NOT NULL DEFAULT 30 CHECK (visit_minutes > 0),
    fee           TEXT NOT NULL DEFAULT '',
    contact       TEXT NOT NULL DEFAULT '',
    tags          TEXT[] NOT NULL DEFAULT '{}',
    images        TEXT[] NOT NULL DEFAULT '{}',
    source        TEXT NOT NULL DEFAULT '',
    published     BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE routes (
    id          BIGSERIAL PRIMARY KEY,
    name        JSONB NOT NULL,
    description JSONB NOT NULL,
    place_ids   BIGINT[] NOT NULL,
    minutes     INT NOT NULL DEFAULT 60,
    published   BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Sheet "Rituals_Events"
CREATE TABLE events (
    id          BIGSERIAL PRIMARY KEY,
    title       JSONB NOT NULL,
    description JSONB NOT NULL,
    preparation JSONB NOT NULL,
    starts_at   TIMESTAMPTZ NOT NULL,
    ends_at     TIMESTAMPTZ NOT NULL,
    place_id    BIGINT REFERENCES places(id) ON DELETE SET NULL,
    contact     TEXT NOT NULL DEFAULT '',
    source      TEXT NOT NULL DEFAULT '',
    published   BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (ends_at >= starts_at)
);
CREATE INDEX events_start_idx ON events(starts_at);

-- Sheet "Shops_Products"
CREATE TABLE shops (
    id          BIGSERIAL PRIMARY KEY,
    name        JSONB NOT NULL,
    kind        TEXT NOT NULL DEFAULT 'shop',
    products    JSONB NOT NULL,
    description JSONB NOT NULL,
    lat         DOUBLE PRECISION,
    lon         DOUBLE PRECISION,
    contact     TEXT NOT NULL DEFAULT '',
    link        TEXT NOT NULL DEFAULT '',
    images      TEXT[] NOT NULL DEFAULT '{}',
    published   BOOLEAN NOT NULL DEFAULT TRUE,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE news (
    id           BIGSERIAL PRIMARY KEY,
    title        JSONB NOT NULL,
    body         JSONB NOT NULL,
    cover_url    TEXT NOT NULL DEFAULT '',
    video_url    TEXT NOT NULL DEFAULT '',
    is_banner    BOOLEAN NOT NULL DEFAULT FALSE,
    published_at TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- RAG knowledge base. Chunks come from published knowledge rows and FAQ rows.
CREATE TABLE kb_chunks (
    id         BIGSERIAL PRIMARY KEY,
    doc_code   TEXT NOT NULL,
    lang       TEXT NOT NULL,
    kind       TEXT NOT NULL DEFAULT '',
    title      TEXT NOT NULL,
    content    TEXT NOT NULL,
    source     TEXT NOT NULL DEFAULT '',
    embedding  vector(256) NOT NULL
);
CREATE INDEX kb_chunks_doc_idx ON kb_chunks(doc_code);
CREATE INDEX kb_chunks_trgm_idx ON kb_chunks USING gin (content gin_trgm_ops);
CREATE INDEX kb_chunks_vec_idx ON kb_chunks USING hnsw (embedding vector_cosine_ops);

-- Sheet "FAQ_AI"
CREATE TABLE faq (
    id        BIGSERIAL PRIMARY KEY,
    category  TEXT NOT NULL,
    question  JSONB NOT NULL,
    answer    JSONB NOT NULL,
    doc_code  TEXT NOT NULL DEFAULT '',
    reviewer  TEXT NOT NULL DEFAULT '',
    published BOOLEAN NOT NULL DEFAULT TRUE
);

-- Chat log. No user identity stored; session_id is a random client id.
CREATE TABLE chat_logs (
    id         BIGSERIAL PRIMARY KEY,
    session_id TEXT NOT NULL,
    lang       TEXT NOT NULL,
    question   TEXT NOT NULL,
    answer     TEXT NOT NULL,
    citations  JSONB NOT NULL DEFAULT '[]',
    refused    BOOLEAN NOT NULL DEFAULT FALSE,
    category   TEXT NOT NULL DEFAULT '',
    latency_ms INT NOT NULL DEFAULT 0,
    helpful    BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX chat_logs_created_idx ON chat_logs(created_at);

-- Sheet "Survey"
CREATE TABLE survey_questions (
    id       BIGSERIAL PRIMARY KEY,
    section  TEXT NOT NULL CHECK (section IN ('satisfaction', 'usability', 'ai_quality')),
    text     JSONB NOT NULL,
    sort     INT NOT NULL DEFAULT 0,
    active   BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE survey_responses (
    id          BIGSERIAL PRIMARY KEY,
    lang        TEXT NOT NULL,
    consent     BOOLEAN NOT NULL CHECK (consent),
    respondent  JSONB NOT NULL DEFAULT '{}',
    answers     JSONB NOT NULL,
    comment     TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE page_views (
    id         BIGSERIAL PRIMARY KEY,
    path       TEXT NOT NULL,
    lang       TEXT NOT NULL,
    visitor    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX page_views_created_idx ON page_views(created_at);

-- Sheet "AI_Test_Set" and experiment results (R1–R9)
CREATE TABLE eval_questions (
    id               BIGSERIAL PRIMARY KEY,
    code             TEXT NOT NULL UNIQUE,
    lang             TEXT NOT NULL,
    question         TEXT NOT NULL,
    reference_answer TEXT NOT NULL,
    relevant_docs    TEXT[] NOT NULL,
    category         TEXT NOT NULL DEFAULT '',
    difficulty       TEXT NOT NULL DEFAULT 'medium',
    reviewer         TEXT NOT NULL DEFAULT ''
);

CREATE TABLE eval_runs (
    id         BIGSERIAL PRIMARY KEY,
    mode       TEXT NOT NULL,
    k          INT NOT NULL,
    note       TEXT NOT NULL DEFAULT '',
    recall_at_k DOUBLE PRECISION,
    mrr        DOUBLE PRECISION,
    n          INT NOT NULL DEFAULT 0,
    status     TEXT NOT NULL DEFAULT 'running',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE eval_results (
    id             BIGSERIAL PRIMARY KEY,
    run_id         BIGINT NOT NULL REFERENCES eval_runs(id) ON DELETE CASCADE,
    question_id    BIGINT NOT NULL REFERENCES eval_questions(id) ON DELETE CASCADE,
    retrieved_docs TEXT[] NOT NULL,
    recall         DOUBLE PRECISION NOT NULL,
    rr             DOUBLE PRECISION NOT NULL,
    answer         TEXT NOT NULL,
    latency_ms     INT NOT NULL,
    expert_correct SMALLINT,
    expert_faithful SMALLINT,
    expert_note    TEXT NOT NULL DEFAULT ''
);
