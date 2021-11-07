CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
DROP TABLE IF EXISTS events CASCADE;
DROP TABLE IF EXISTS snapshots CASCADE;

CREATE TABLE IF NOT EXISTS snapshots
(
    snapshot_id    UUID PRIMARY KEY         DEFAULT uuid_generate_v4(),
    aggregate_id   VARCHAR(250) UNIQUE NOT NULL CHECK ( aggregate_id <> '' ),
    aggregate_type VARCHAR(250)        NOT NULL CHECK ( aggregate_type <> '' ),
    data           BYTEA,
    metadata       BYTEA,
    version        SERIAL              NOT NULL,
    timestamp      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (aggregate_id)
);

CREATE INDEX IF NOT EXISTS aggregate_id_aggregate_type_idx ON snapshots USING btree (aggregate_id, version);