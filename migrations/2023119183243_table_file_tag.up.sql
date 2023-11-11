CREATE TABLE IF NOT EXISTS file_tag(
    id SERIAL PRIMARY KEY,
    tag_id BIGINT NOT NULL,
    file_object_id BIGINT NOT NULL,
    CONSTRAINT file_id_tag_id_uq UNIQUE (tag_id,file_object_id),
    CONSTRAINT file_tag_tag_id_fk FOREIGN KEY (tag_id) REFERENCES tag (id) ON DELETE CASCADE,
    CONSTRAINT file_tag_file_object_id_fk FOREIGN KEY (file_object_id) REFERENCES file_object (id) ON DELETE CASCADE
);

