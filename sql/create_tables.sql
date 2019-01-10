CREATE TABLE documents (
    id                  VARCHAR(100),
    type                VARCHAR(100),
    version             INTEGER NOT NULL,
    time_created        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    time_updated        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    content             JSONB,
    PRIMARY KEY         (id, type)
);