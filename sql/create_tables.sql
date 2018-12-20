CREATE TABLE orders (
    id                  uuid PRIMARY KEY,
    version             integer NOT NULL,
    status              varchar(50) NOT NULL,
    time_created        timestamp DEFAULT CURRENT_TIMESTAMP,
    time_updated        timestamp DEFAULT CURRENT_TIMESTAMP,
    time_placed         timestamp NOT NULL,
    details             jsonb
);

CREATE TABLE process_defs (
    id                  VARCHAR(100) PRIMARY KEY,
    description         text,
    time_created        timestamp DEFAULT CURRENT_TIMESTAMP,
    time_updated        timestamp DEFAULT CURRENT_TIMESTAMP,
    details             jsonb
)