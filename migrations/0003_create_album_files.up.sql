CREATE TABLE album_files
(
    id         INTEGER  NOT NULL PRIMARY KEY,
    album_id   INTEGER  NOT NULL REFERENCES albums (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    file_id    INTEGER  NOT NULL REFERENCES files (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now', 'utc')),
    updated_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now', 'utc'))
);

-- CREATE TRIGGER album_files_updated_at
--     AFTER UPDATE
--     ON album_files
--     FOR EACH ROW
-- BEGIN
--     UPDATE album_files
--     SET updated_at = strftime('%Y-%m-%d %H:%M:%f', 'now', 'utc')
--     WHERE id = old.id;
-- END;

CREATE INDEX album_files_album_id ON album_files (album_id);
CREATE INDEX album_files_file_id ON album_files (file_id);
