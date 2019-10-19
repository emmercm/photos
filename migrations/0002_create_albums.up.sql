CREATE TABLE albums
(
    id         INTEGER  NOT NULL PRIMARY KEY,
    path       TEXT UNIQUE,
    title      TEXT     NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now', 'utc')),
    updated_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now', 'utc')),
    deleted_at DATETIME
);

-- CREATE TRIGGER albums_updated_at
--     AFTER UPDATE
--     ON albums
--     FOR EACH ROW
-- BEGIN
--     UPDATE albums
--     SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now', 'utc')
--     WHERE id = old.id;
-- END;
