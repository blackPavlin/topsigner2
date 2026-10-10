BEGIN;
CREATE TYPE user_role AS ENUM ('USER', 'ADMIN');

CREATE TABLE users (
    id 		        BIGINT 	    GENERATED ALWAYS AS IDENTITY,
    vk_user_id      BIGINT      DEFAULT NULL,
    email           TEXT        DEFAULT NULL,
    password_hash   TEXT        DEFAULT NULL,
    role            user_role   NOT NULL DEFAULT 'USER',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    CONSTRAINT users_email_unique UNIQUE (email),
    CONSTRAINT users_vk_user_id_unique UNIQUE (vk_user_id),
    CONSTRAINT users_identity_check CHECK (email IS NOT NULL OR vk_user_id IS NOT NULL)
);

CREATE TYPE session_auth_type AS ENUM ('PASSWORD', 'VK_OAUTH');

CREATE TABLE sessions (
    id                 UUID              NOT NULL DEFAULT uuidv7(),
    user_id            BIGINT            NOT NULL,
    auth_type          session_auth_type NOT NULL,
    device_id          TEXT              NOT NULL,
    ip                 INET              NOT NULL,
    user_agent         TEXT              NOT NULL,
    refresh_token_hash TEXT              NOT NULL,
    expires_at         TIMESTAMPTZ       NOT NULL,
    created_at         TIMESTAMPTZ       NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ       NOT NULL DEFAULT now(),
 
    PRIMARY KEY (id),
    CONSTRAINT sessions_refresh_token_hash_unique UNIQUE (refresh_token_hash),
    CONSTRAINT sessions_user_id_auth_type_device_id_unique UNIQUE (user_id, auth_type, device_id),
    CONSTRAINT sessions_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX sessions_user_id_idx ON sessions (user_id);
CREATE INDEX sessions_expires_at_idx ON sessions (expires_at);

CREATE TABLE vk_sessions (
    id                UUID        NOT NULL DEFAULT uuidv7(),
    session_id        UUID        NOT NULL,
    device_id         TEXT        NOT NULL,
    access_token_enc  TEXT        NOT NULL,
    refresh_token_enc TEXT        NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
 
    PRIMARY KEY (id),
    CONSTRAINT vk_sessions_session_id_fk FOREIGN KEY (session_id) REFERENCES sessions (id) ON DELETE CASCADE
);

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

CREATE TABLE groups (
    id 		         BIGINT 	 GENERATED ALWAYS AS IDENTITY,
    user_id          BIGINT      NOT NULL,
    external_id      BIGINT      NOT NULL,
    access_token_enc TEXT        NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    CONSTRAINT groups_user_id_external_id_unique UNIQUE (user_id, external_id),
    CONSTRAINT groups_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
COMMIT;
