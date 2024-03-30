-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS users
(
    id             CHAR(27)     NOT NULL,
    email          VARCHAR(255) NOT NULL,
    password       VARCHAR(255) NOT NULL,
    is_admin       TINYINT(1)   NULL,
    created_at     TIMESTAMP    NOT NULL,
    updated_at     TIMESTAMP    NULL,
    deleted_at     TIMESTAMP    NULL,
    PRIMARY KEY (ID)
);

INSERT INTO users (id,email,password,is_admin,created_at) VALUES("200KPME8UZ4tjpP0IqMyyAizmsy", "system@bluffy.de","$2a$14$5sPvPBjLbUf4MOTcu1.izOhdTOx..AAcmGprzHkCNy2ckdRiGadxe", 1, NOW());
COMMIT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS users;