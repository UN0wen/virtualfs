DROP TABLE IF EXISTS closure;

DROP TABLE IF EXISTS filedirs;

DROP TYPE IF EXISTS filetype;

CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TABLE IF NOT EXISTS filedirs (
    id serial PRIMARY KEY,
    parentdirid int REFERENCES filedirs (id),
    name text NOT NULL,
    data text NOT NULL DEFAULT '',
    created timestamptz NOT NULL DEFAULT now(),
    updated timestamptz NOT NULL DEFAULT now(),
    size int DEFAULT 0,
    type text NOT NULL,
    path ltree UNIQUE
);

CREATE INDEX path_gist_idx ON filedirs USING GIST (path);

CREATE INDEX path_btree_idx ON filedirs USING btree (path);

-------------------------------------------
-------------------------------------------
-- Triggers to manage sizes automatically--
-------------------------------------------
-------------------------------------------
CREATE OR REPLACE FUNCTION insert_size ()
    RETURNS TRIGGER
    AS $$
BEGIN
    IF NEW.type = 'file' THEN
        UPDATE
            filedirs AS f
        SET
            size = f2.size + NEW.size
        FROM (
            SELECT
                path,
                size
            FROM
                filedirs
            WHERE
                path @> NEW.path
                AND TYPE = 'directory') AS f2 (path,
            size)
    WHERE
        f2.path = f.path
            AND f.type = 'directory';
    END IF;
    RETURN NEW;
END;
$$
LANGUAGE PLPGSQL;

CREATE TRIGGER insert_size
    AFTER INSERT ON filedirs
    FOR EACH ROW
    EXECUTE PROCEDURE insert_size ();

CREATE OR REPLACE FUNCTION update_size ()
    RETURNS TRIGGER
    AS $$
BEGIN
    IF NEW.type = 'file' THEN
        UPDATE
            filedirs AS f
        SET
            size = f2.size + (NEW.size - OLD.size)
        FROM (
            SELECT
                path,
                size
            FROM
                filedirs
            WHERE
                path @> NEW.path
                AND TYPE = 'directory') AS f2 (path,
            size)
    WHERE
        f2.path = f.path
            AND f.type = 'directory';
    END IF;
    RETURN NEW;
END;
$$
LANGUAGE PLPGSQL;

CREATE TRIGGER update_size
    AFTER UPDATE ON filedirs
    FOR EACH ROW
    EXECUTE PROCEDURE update_size ();

CREATE OR REPLACE FUNCTION delete_size ()
    RETURNS TRIGGER
    AS $$
BEGIN
    IF OLD.type = 'file' THEN
        UPDATE
            filedirs AS f
        SET
            size = f2.size - OLD.size
        FROM (
            SELECT
                path,
                size
            FROM
                filedirs
            WHERE
                path @> OLD.path
                AND TYPE = 'directory') AS f2 (path,
            size)
    WHERE
        f2.path = f.path
            AND f.type = 'directory';
    END IF;
    RETURN OLD;
END;
$$
LANGUAGE PLPGSQL;

CREATE TRIGGER delete_size
    BEFORE DELETE ON filedirs
    FOR EACH ROW
    EXECUTE PROCEDURE delete_size ();

------------------
------------------
-- End triggers --
------------------
------------------


-- Root dir
INSERT INTO filedirs (id, parentdirid, name, TYPE, path)
    VALUES (DEFAULT, lastval(), '/', 'directory', 'root');

-- Test values
INSERT INTO filedirs (id, parentdirid, name, TYPE, path)
    VALUES (DEFAULT, 1, 'test', 'directory', 'root.test');

INSERT INTO filedirs (id, parentdirid, name, TYPE, path)
    VALUES (DEFAULT, 2, 'test2', 'directory', 'root.test.test2');

INSERT INTO filedirs (id, parentdirid, name, TYPE, path, size)
    VALUES (DEFAULT, 1, 'file1', 'file', 'root.file1', 5);

INSERT INTO filedirs (id, parentdirid, name, TYPE, path, size)
    VALUES (DEFAULT, 1, 'file2', 'file', 'root.file2', 10);

INSERT INTO filedirs (id, parentdirid, name, TYPE, path, size)
    VALUES (DEFAULT, 2, 'file3', 'file', 'root.test.file3', 3);

INSERT INTO filedirs (id, parentdirid, name, TYPE, path, size)
    VALUES (DEFAULT, 3, 'file4', 'file', 'root.test.test2.file4', 8);


