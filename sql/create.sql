CREATE TABLE IF NOT EXISTS activities
(
    name text NULL,
    postcode text,
    sunny boolean,
    CONSTRAINT name PRIMARY KEY (name)
    );

ALTER TABLE IF EXISTS activities
    OWNER to activities;
