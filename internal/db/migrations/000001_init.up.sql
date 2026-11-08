CREATE TYPE packing_type AS ENUM(
    'Quantity',
    'Volume',
    'Weight'
);

CREATE TABLE "products"(
    "packing_type" packing_type NOT NULL DEFAULT ('Quantity'),
    "packing_coefficient_types" NUMERIC[] NOT NULL DEFAULT ARRAY[1],
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "name" varchar UNIQUE NOT NULL,
    "description" text,
    "category" uuid,
    "image_url" text,
    "slug" varchar,
    "price" NUMERIC  DEFAULT 0,
    "count" NUMERIC DEFAULT 0,
    "discount" varchar(4),
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "users"(
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "email" varchar UNIQUE,
    "phone_number" varchar(20) UNIQUE,
    "expired_time" timestamptz,
    "password" text NOT NULL,
    "token" text,
    "is_active" boolean DEFAULT FALSE,
    "is_admin" boolean DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    CONSTRAINT "at least one contact method" CHECK (email IS NOT NULL OR phone_number IS NOT NULL)
);

CREATE TABLE "categories"(
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "name" varchar UNIQUE NOT NULL,
    "description" text,
    "parent" uuid,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "files"(
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "name" varchar UNIQUE NOT NULL,
    "url" varchar,
    "product" uuid,
    "category" uuid,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users"("email");

CREATE INDEX ON "users"("phone_number");

ALTER TABLE "categories"
    ADD FOREIGN KEY ("parent") REFERENCES "categories"("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "products"
    ADD FOREIGN KEY ("category") REFERENCES "categories"("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "files"
    ADD FOREIGN KEY ("product") REFERENCES "products"("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "files"
    ADD FOREIGN KEY ("category") REFERENCES "categories"("id") DEFERRABLE INITIALLY IMMEDIATE;

