CREATE TABLE "users"(
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "email" varchar UNIQUE,
    "phone_number" varchar(20) UNIQUE,
    "expired_time" timestamptz,
    "password" text NOT NULL,
    "token" text,
    "is_active" boolean,
    "is_admin" boolean,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    CONSTRAINT at_least_one_contact_method CHECK (email IS NOT NULL OR phone_number IS NOT NULL)
);

CREATE TABLE "categories"(
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "name" varchar UNIQUE NOT NULL,
    "description" text,
    "parent" uuid,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users"("email");

CREATE INDEX ON "users"("phone_number");

ALTER TABLE "categories"
    ADD FOREIGN KEY ("parent") REFERENCES "categories"("id") DEFERRABLE INITIALLY IMMEDIATE;

