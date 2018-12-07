CREATE TABLE orders (  
  id uuid PRIMARY KEY,
  version integer NOT NULL,
  time_created timestamp NOT NULL,
  time_placed timestamp NOT NULL,
  time_updated timestamp default current_timestamp,
  data jsonb
);