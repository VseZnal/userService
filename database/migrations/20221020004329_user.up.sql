BEGIN;

CREATE OR REPLACE FUNCTION pseudo_encrypt(value BIGINT) returns BIGINT AS $$
DECLARE
    l1  BIGINT;
    l2  BIGINT;
    r1  BIGINT;
    r2  BIGINT;
    i   BIGINT := 0;
BEGIN
    l1:= (value >> 16) & 65535;
    r1:= value & 65535;
    WHILE i < 3 LOOP
            l2 := r1;
            r2 := l1 # ((((1366 * r1 + 150889) % 714025) / 714025.0) * 32767)::BIGINT;
            l1 := l2;
            r1 := r2;
            i := i + 1;
        END LOOP;
    return ((r1 << 16) + l1);
END;
$$ LANGUAGE plpgsql strict immutable;

CREATE SEQUENCE users_id_seq;

CREATE TABLE users (
    id          BIGINT PRIMARY KEY NOT NULL DEFAULT pseudo_encrypt(nextval('users_id_seq')::BIGINT),
    username    TEXT NOT NULL ,
    password    TEXT NOT NULL
);

COMMIT;