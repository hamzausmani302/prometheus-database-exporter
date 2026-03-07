-- =============================================================================
-- Test schema for database-prometheus-exporter E2E tests
-- Domain: oil & gas drilling operations
-- Run against: exporter database (docker-compose postgres on :5433)
-- =============================================================================

-- Cleanup (safe to re-run)
DROP TABLE IF EXISTS t_ddr       CASCADE;
DROP TABLE IF EXISTS t_operation CASCADE;
DROP TABLE IF EXISTS t_well      CASCADE;
DROP TABLE IF EXISTS t_rig       CASCADE;

-- =============================================================================
-- Tables
-- =============================================================================

-- Rigs (drilling equipment assigned to wells)
CREATE TABLE t_rig (
    id_t_rig  SERIAL PRIMARY KEY,
    rig_name  VARCHAR(100) NOT NULL
);

-- Wells (oil & gas wells being monitored)
CREATE TABLE t_well (
    id_t_well      SERIAL PRIMARY KEY,
    well_name      VARCHAR(100) NOT NULL,
    well_id        VARCHAR(50)  NOT NULL,  -- external asset identifier
    is_live        BOOLEAN      NOT NULL DEFAULT false,  -- currently drilling
    to_load        BOOLEAN      NOT NULL DEFAULT true,   -- included in metric exports
    ts_last_mudlog TIMESTAMP                             -- last sensor data received
);

-- Operations: which rig is (or was) on a well, and when it started.
-- The most-recent row per well (by start_date) is the active assignment.
CREATE TABLE t_operation (
    id_t_operation SERIAL PRIMARY KEY,
    id_t_well      INT NOT NULL REFERENCES t_well(id_t_well),
    id_t_rig       INT NOT NULL REFERENCES t_rig(id_t_rig),
    start_date     TIMESTAMP NOT NULL
);

-- Daily drilling report: one row per activity period, capturing its duration.
CREATE TABLE t_ddr (
    id_t_ddr      SERIAL PRIMARY KEY,
    id_t_well     INT     NOT NULL REFERENCES t_well(id_t_well),
    duration_sec  NUMERIC NOT NULL,  -- activity duration in seconds
    ts_start_utc  TIMESTAMP NOT NULL,
    ts_end_utc    TIMESTAMP NOT NULL
);

-- =============================================================================
-- Seed data
-- =============================================================================

INSERT INTO t_rig (rig_name) VALUES
    ('Rig Alpha'),   -- id 1
    ('Rig Beta'),    -- id 2
    ('Rig Gamma');   -- id 3

-- 5 wells with different live/load combinations:
--   is_live=true  + to_load=true  → included in live queries  (wells 1-3)
--   is_live=false + to_load=true  → offline, counted by status (well 4)
--   is_live=false + to_load=false → fully excluded             (well 5)
INSERT INTO t_well (well_name, well_id, is_live, to_load, ts_last_mudlog) VALUES
    ('Permian A-01',    'A-01', true,  true,  NOW() - INTERVAL '5 minutes'),   -- id 1
    ('Permian A-02',    'A-02', true,  true,  NOW() - INTERVAL '2 minutes'),   -- id 2
    ('Eagle Ford B-01', 'B-01', true,  true,  NOW() - INTERVAL '10 minutes'),  -- id 3
    ('Eagle Ford C-01', 'C-01', false, true,  NOW() - INTERVAL '1 day'),       -- id 4
    ('Bakken D-01',     'D-01', false, false, NOW() - INTERVAL '3 days');      -- id 5

-- Operations — each well may have multiple rows; most-recent determines active rig.
-- Well 1: currently Rig Alpha (2025-01), previously Rig Beta (2024-06)
INSERT INTO t_operation (id_t_well, id_t_rig, start_date) VALUES
    (1, 2, '2024-06-01 00:00:00'),
    (1, 1, '2025-01-01 00:00:00');

-- Well 2: currently Rig Alpha
INSERT INTO t_operation (id_t_well, id_t_rig, start_date) VALUES
    (2, 1, '2025-02-01 00:00:00');

-- Well 3: currently Rig Beta
INSERT INTO t_operation (id_t_well, id_t_rig, start_date) VALUES
    (3, 2, '2025-01-15 00:00:00');

-- Well 4: currently Rig Gamma (offline — verifies is_live filter works)
INSERT INTO t_operation (id_t_well, id_t_rig, start_date) VALUES
    (4, 3, '2024-12-01 00:00:00');

-- DDR records (duration data used by pipeline and simple duration queries)
-- Well 1 (Permian A-01):    1200 + 2400 = 3600 sec total
INSERT INTO t_ddr (id_t_well, duration_sec, ts_start_utc, ts_end_utc) VALUES
    (1, 1200, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'  + INTERVAL '1200 seconds'),
    (1, 2400, NOW() - INTERVAL '1 day',  NOW() - INTERVAL '1 day'   + INTERVAL '2400 seconds');

-- Well 2 (Permian A-02):    3600 + 3600 = 7200 sec total
INSERT INTO t_ddr (id_t_well, duration_sec, ts_start_utc, ts_end_utc) VALUES
    (2, 3600, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'  + INTERVAL '3600 seconds'),
    (2, 3600, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'  + INTERVAL '3600 seconds');

-- Well 3 (Eagle Ford B-01): 1800 sec total
INSERT INTO t_ddr (id_t_well, duration_sec, ts_start_utc, ts_end_utc) VALUES
    (3, 1800, NOW() - INTERVAL '1 day',  NOW() - INTERVAL '1 day'   + INTERVAL '1800 seconds');
