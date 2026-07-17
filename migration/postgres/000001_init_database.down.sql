BEGIN;
DROP TABLE fonts;
DROP TABLE images;
DROP INDEX sessions_user_id_idx;
DROP TABLE sessions;
DROP INDEX images_user_id_idx;
DROP TABLE users;
DROP TYPE user_role;
COMMIT;
