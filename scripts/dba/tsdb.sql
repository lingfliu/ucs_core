create table if not exist device_meta (
    id SERIAL PRIMARY KEY,
    mac VARCHAR NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE(device_id, key)
);