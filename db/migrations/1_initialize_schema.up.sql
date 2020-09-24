BEGIN;
CREATE TABLE "users"
(
    "id"       bigserial,
    "login"    text NOT NULL UNIQUE,
    "password" text NOT NULL,
    PRIMARY KEY ("id"),
    UNIQUE ("login")
);

CREATE TABLE "bikes"
(
    "id"   bigserial,
    "name" text NOT NULL,
    "lat"  real NOT NULL,
    "lon"  real NOT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "rentals"
(
    "id"        bigserial,
    "bike_id"   bigint  NOT NULL,
    "user_id"   bigint  NOT NULL,
    "completed" boolean NOT NULL DEFAULT false,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("bike_id") REFERENCES "bikes" ("id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
COMMIT;
