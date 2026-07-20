BEGIN;
CREATE TYPE user_role AS ENUM ('USER', 'ADMIN');

CREATE TABLE users (
    id 		   BIGINT 	   GENERATED ALWAYS AS IDENTITY,
    email      TEXT        NOT NULL,
    role       user_role   NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id)
);

CREATE TABLE sessions (
    id                 UUID        NOT NULL DEFAULT uuidv7(),
    user_id            BIGINT      NOT NULL,
    ip                 TEXT        NOT NULL,
    refresh_token_hash TEXT        NOT NULL,
    expires_at         TIMESTAMPTZ NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    CONSTRAINT sessions_refresh_token_hash_unique UNIQUE (refresh_token_hash),
    CONSTRAINT sessions_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);

CREATE TABLE fonts (
    id 		   BIGINT 	   GENERATED ALWAYS AS IDENTITY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    CONSTRAINT fonts_name_unique UNIQUE (name)
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
