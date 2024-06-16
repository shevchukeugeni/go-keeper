CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id         UUID         NOT NULL PRIMARY KEY,
    login      varchar(128) NOT NULL,
    password   text         NOT NULL,

    created_at TIMESTAMP    NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS users_login_idx ON users (login);

CREATE TABLE IF NOT EXISTS texts
(
    user_id     uuid,
    id          uuid,
    key         varchar   NOT NULL,
    data        text,
    metadata    varchar,
    uploaded_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
        DEFERRABLE INITIALLY DEFERRED
);

CREATE UNIQUE INDEX IF NOT EXISTS texts_idx ON texts (id);

CREATE TABLE IF NOT EXISTS cards
(
    user_id     uuid,
    id          uuid,
    card        varchar   NOT NULL,
    expiration  varchar,
    cvv         varchar,
    metadata    varchar,
    uploaded_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
        DEFERRABLE INITIALLY DEFERRED
);

CREATE UNIQUE INDEX IF NOT EXISTS cards_idx ON cards (id);

CREATE TABLE IF NOT EXISTS credentials
(
    user_id     uuid,
    id          uuid,
    site        varchar   NOT NULL,
    login       varchar,
    password    varchar,
    metadata    varchar,
    uploaded_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
        DEFERRABLE INITIALLY DEFERRED
);

CREATE UNIQUE INDEX IF NOT EXISTS creds_idx ON credentials (id);