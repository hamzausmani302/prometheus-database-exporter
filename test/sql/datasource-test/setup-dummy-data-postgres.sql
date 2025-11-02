CREATE TYPE task_status AS ENUM ('pending', 'success', 'failed', 'cancelled');
CREATE TABLE products (
    task_id SERIAL PRIMARY KEY,
    task_id_ext VARCHAR(100),
    task_query VARCHAR(200),
    task_name VARCHAR(100) NOT NULL,
    run_date TIMESTAMP,
    created_date DATE DEFAULT CURRENT_DATE,
    task_status task_status NOT NULL DEFAULT 'pending'
);

INSERT INTO products (task_name,task_id_ext, task_query, run_date, task_status) VALUES ('get-failed-task', '1','select 1', '2025-03-03T00:00:00', 'pending' );
INSERT INTO products (task_name,task_id_ext, task_query, run_date, task_status) VALUES ('get-failed-task', '2','select 2', '2025-03-04T00:00:00', 'failed' );
INSERT INTO products (task_name,task_id_ext, task_query, run_date, task_status) VALUES ('get-failed-task', '2','select 2', '2025-03-04T00:00:00', 'cancelled' );

