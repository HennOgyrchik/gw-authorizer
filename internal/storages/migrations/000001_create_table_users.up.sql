CREATE TABLE IF NOT EXISTS users (
                              id serial4 NOT NULL,
                              login varchar(10) NOT NULL,
                              "password" text NOT NULL,
                              email varchar(30) NOT NULL,
                              token text,
                              CONSTRAINT users_email_key UNIQUE (email),
                              CONSTRAINT users_id_key UNIQUE (id),
                              CONSTRAINT users_login_key UNIQUE (login)
);