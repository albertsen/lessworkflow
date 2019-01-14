CREATE TABLE documents (
    id                  VARCHAR(100),
    type                VARCHAR(100),
    version             INTEGER NOT NULL,
    time_created        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    time_updated        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    content             JSONB,
    PRIMARY KEY         (id, type)
);

CREATE TABLE processes (
    id                  UUID,
    time_created        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    time_updated        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status              VARCHAR(100),
    document_url        VARCHAR(1000) NOT NULL,
    process_def_url     VARCHAR(1000) NOT NULL,
    PRIMARY KEY         (id)
);