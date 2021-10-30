CREATE TABLE IF NOT EXISTS shows
(
    id          serial PRIMARY KEY,
    title       VARCHAR(100) UNIQUE,
    weekday     VARCHAR(3),
    timeslot    VARCHAR(11),
    description TEXT,
    created_at  TIMESTAMPTZ,
    updated_at  TIMESTAMPTZ
)