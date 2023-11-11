BEGIN;
CREATE TABLE IF NOT EXISTS tag(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE UNIQUE INDEX tag_name_index_uq on tag(name) INCLUDE (id);

COMMIT;
