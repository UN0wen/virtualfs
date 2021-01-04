
DROP TABLE IF EXISTS filedirs;

CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TABLE IF NOT EXISTS filedirs (
    id serial PRIMARY KEY,
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
            size = f2.size + NEW.size,
            updated = now()
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
        IF NEW.path = OLD.path THEN
            UPDATE
                filedirs AS f
            SET
                size = f2.size + (NEW.size - OLD.size),
                updated = now()
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
        ELSE
            -- Delete from old
            UPDATE
                filedirs AS f
            SET
                size = f2.size - OLD.size,
                updated = now()
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
            -- Add to new
            UPDATE
                filedirs AS f
            SET
                size = f2.size + NEW.size,
                updated = now()
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
            size = f2.size - OLD.size,
            updated = now()
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
INSERT INTO filedirs (id, name, TYPE, path)
    VALUES (DEFAULT, '/', 'directory', 'root');