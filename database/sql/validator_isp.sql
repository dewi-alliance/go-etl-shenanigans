DROP TABLE IF EXISTS "public"."validator_isp";

CREATE TABLE "public"."validator_isp" (
    "address" text NOT NULL,
    "isp" text,
    "geo_data" text,
    PRIMARY KEY ("address")
);

