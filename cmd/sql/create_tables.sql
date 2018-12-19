CREATE TYPE monetary_amount AS (
    value               bigint,
    currency            char(3)
);

CREATE TYPE line_item AS (
    product_id          varchar(100),
    product_description text,
    count               integer,
    item_price          monetary_amount,
    total_price         monetary_amount,
    properties          jsonb
);

CREATE TABLE orders (
    id                  uuid PRIMARY KEY,
    version             integer NOT NULL,
    status              varchar(50) NOT NULL,
    total_price         monetary_amount NOT NULL,
    time_created        timestamp DEFAULT CURRENT_TIMESTAMP,
    time_updated        timestamp DEFAULT CURRENT_TIMESTAMP,
    time_placed         timestamp NOT NULL,
    line_items          line_item[] NOT NULL,
    properties          jsonb
);