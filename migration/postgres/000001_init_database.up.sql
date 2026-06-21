BEGIN;
CREATE TABLE images (
    id 		   BIGINT 	   GENERATED ALWAYS AS IDENTITY,
    name       TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (id),
    CONSTRAINT image_name_unique UNIQUE (name)
);
COMMIT;
