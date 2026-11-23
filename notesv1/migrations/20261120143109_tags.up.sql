CREATE TABLE tags (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL,
    name       TEXT        NOT NULL,

    CONSTRAINT tags_user_id_name_key UNIQUE (user_id, name)
);

CREATE TABLE note_tags (
    note_id BIGINT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    tag_id  UUID NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,

    CONSTRAINT note_tags_pkey PRIMARY KEY (note_id, tag_id)
);