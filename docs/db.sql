CREATE TABLE test_timestamp (
  id            SERIAL PRIMARY KEY,
  create_time   BIGINT NOT NULL
);

CREATE TABLE test_datetime (
  id            SERIAL PRIMARY KEY,
  create_time   TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL DEFAULT now()
);