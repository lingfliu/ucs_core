create table if not exists device_meta (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    mac VARCHAR NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    UNIQUE(device_id, key)
);

create table if not exists project (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR NOT NULL,
    description TEXT NOT NULL,
    UNIQUE(name)
);