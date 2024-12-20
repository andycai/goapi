--
-- File generated with SQLiteStudio v3.4.4 on Thu Dec 12 19:55:20 2024
--
-- Text encoding used: UTF-8
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- Table: task_logs
DROP TABLE IF EXISTS task_logs;

CREATE TABLE IF NOT EXISTS task_logs (
    id         INTEGER      PRIMARY KEY AUTOINCREMENT,
    task_id    INTEGER      NOT NULL,
    status     VARCHAR (20) NOT NULL,
    output     TEXT,
    error      TEXT,
    start_time DATETIME     NOT NULL,
    end_time   DATETIME,
    duration   INTEGER,
    created_at DATETIME     DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     DEFAULT (CURRENT_TIMESTAMP),
    FOREIGN KEY (
        task_id
    )
    REFERENCES tasks (id) ON DELETE CASCADE
);


COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
