
CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    version integer NOT NULL
);

INSERT INTO schema_migrations (version) VALUES (0);


CREATE TABLE IF NOT EXISTS ballots (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

