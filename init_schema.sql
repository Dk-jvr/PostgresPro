DROP TABLE IF EXISTS Scripts CASCADE;
DROP TABLE IF EXISTS OutPutCommands;

CREATE TABLE IF NOT EXISTS Scripts (
    script_id SERIAL PRIMARY KEY,
    script TEXT NOT NULL,
    script_type TEXT NOT NULL
);

INSERT INTO Scripts(script, script_type) VALUES ('#!/bin/bash\n\n for ((i = 1; i <= 10; i++)); do\n    echo \"$(date +''%H:%M:%S'') - Some data $i\"\n    sleep 1\n done', '/bin/bash');

CREATE TABLE IF NOT EXISTS OutPutScripts (
    script_id INTEGER NOT NULL REFERENCES Scripts(script_id),
    script TEXT NOT NULL,
    script_type TEXT,
    output TEXT,
    output_time TIMESTAMP
)