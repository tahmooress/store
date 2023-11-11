BEGIN;
CREATE TABLE IF NOT EXISTS total_storage(
    size BIGINT DEFAULT 0
);

INSERT INTO total_storage(size) VALUES(0);
COMMIT;
