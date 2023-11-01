CREATE TABLE users (
    id uuid PRIMARY KEY NOT NULL ,
    created TIMESTAMPTZ NOT NULL DEFAULT now(),
    email TEXT NOT NULL UNIQUE,
    admin bool NOT NULL DEFAULT false,
    hashed_password TEXT NOT NULL
);
