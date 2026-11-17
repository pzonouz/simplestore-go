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
    "is_service" boolean DEFAULT False,
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
    "slug" text,
    "priority" NUMERIC DEFAULT 100,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "files"(
    "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
    "name" varchar UNIQUE NOT NULL,
    "url" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "products_files" (
    product_id uuid NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    file_id  uuid NOT NULL REFERENCES files(id)  ON DELETE CASCADE,
    PRIMARY KEY (product_id, file_id)
);



CREATE INDEX ON "users"("email");
CREATE INDEX ON "users"("phone_number");

CREATE INDEX ON "categoires"("slug");
CREATE INDEX ON "products"("slug");


ALTER TABLE "categories"
ADD FOREIGN KEY ("parent") REFERENCES "categories"("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "products"
ADD FOREIGN KEY ("category") REFERENCES "categories"("id") DEFERRABLE INITIALLY IMMEDIATE;

CREATE OR REPLACE FUNCTION category_tree(p_parent uuid DEFAULT NULL)
RETURNS jsonb
LANGUAGE sql
STABLE
AS $$
    SELECT COALESCE(jsonb_agg(node ORDER BY c.priority ASC NULLS LAST), '[]'::jsonb)
    FROM (
        SELECT jsonb_build_object(
                   'id',       c.id,
                   'name',     c.name,
                   'slug',     c.slug,
                   'priority', c.priority,
                   'children', category_tree(c.id)      -- <-- recursive call
               ) AS node,
               c.priority,
               c.name
        FROM categories c
        WHERE c.parent IS NOT DISTINCT FROM p_parent   -- NULL-safe for roots
    ) s;
$$;

CREATE OR REPLACE VIEW category_tree_view AS
SELECT category_tree() AS tree;
