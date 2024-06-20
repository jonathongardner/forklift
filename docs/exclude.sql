CREATE TABLE files (id INT PRIMARY KEY, path TEXT NOT NULL, type TEXT NOT NULL, parent_id INT NOT NULL);
CREATE INDEX idx_files_parent_id ON files (parent_id);
DELETE FROM files;
INSERT INTO files (id, path, type, parent_id) VALUES
    (1, 'ubuntu.iso', 'mbr', 0),
    (2, '/bin/cool', 'binary', 1),
    (3, '/bin/yo.tar', 'tar', 1),
    (4, '/bin/yo.tar/coolo', 'binary', 3),
    (5, '/bin/yo.tar/yoho', 'binary', 3),
    (6, '/bin/yo.tar/another/tar', 'tar', 3),
    (7, '/bin/yo.tar/another/tar/some-file', 'ascii', 6),
    (8, '/bin/coolo', 'binary', 1),
    (9, '/bin/yo.zip', 'zip', 1),
    (10, '/bin/yo.zip/yeah', 'binary', 9);

WITH RECURSIVE to_unpack (id, path, type, children) AS (
    SELECT 0, '', '', TRUE

    UNION

    SELECT f.id, f.path, f.type, (SELECT 1 FROM files f2 WHERE f.type NOT IN ('tar', 'mbr') AND f2.parent_id = f.id LIMIT 1)
    FROM files f
    JOIN to_unpack t ON t.id = f.parent_id
    WHERE t.type NOT IN ('tar', 'mbr')
)

SELECT * FROM to_unpack WHERE children IS NULL;