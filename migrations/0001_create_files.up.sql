CREATE TABLE files
(
    id         INTEGER  NOT NULL PRIMARY KEY,
    path       TEXT     NOT NULL UNIQUE,
    fast_hash  TEXT,
    slow_hash  TEXT,
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now', 'utc')),
    updated_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now', 'utc')),
    deleted_at DATETIME
);

-- CREATE TRIGGER files_updated_at
--     AFTER UPDATE
--     ON files
--     FOR EACH ROW
-- BEGIN
--     UPDATE files
--     SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now', 'utc')
--     WHERE path = old.path;
-- END;
