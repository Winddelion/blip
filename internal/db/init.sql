CREATE TABLE IF NOT EXISTS checks (
    id  BIGSERIAL PRIMARY KEY,
    url TEXT    NOT NULL,
    status_code INTEGER,
    response_time INTEGER,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_checkls_url_time
    ON checks (url, checked_at DESC);