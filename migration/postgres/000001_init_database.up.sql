BEGIN;
CREATE TYPE user_role AS ENUM ('USER', 'ADMIN');

CREATE TABLE users (
    id 		   BIGINT 	   GENERATED ALWAYS AS IDENTITY,
    role       user_role   NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id)
);

CREATE TABLE images (
    id 		   BIGINT 	   GENERATED ALWAYS AS IDENTITY,
    user_id    BIGINT      NOT NULL,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    CONSTRAINT images_name_unique UNIQUE (name),
    CONSTRAINT images_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX images_user_id_idx ON images (user_id);
COMMIT;
