--
-- File generated with SQLiteStudio v3.4.4 on Thu Dec 12 19:54:49 2024
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: tasks
DROP TABLE IF EXISTS tasks;

CREATE TABLE IF NOT EXISTS tasks (
    id          INTEGER       PRIMARY KEY AUTOINCREMENT,
    name        VARCHAR (100) NOT NULL,
    description TEXT,
    type        VARCHAR (20)  NOT NULL,
    script      TEXT,
    url         VARCHAR (255),
    method      VARCHAR (10),
    headers     TEXT,
    body        TEXT,
    timeout     INTEGER       DEFAULT 300,
    status      VARCHAR (20)  DEFAULT 'inactive',
    created_at  DATETIME      DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME      DEFAULT CURRENT_TIMESTAMP,
    cron_expr   VARCHAR (100),
    enable_cron INTEGER       DEFAULT (0) 
);


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
