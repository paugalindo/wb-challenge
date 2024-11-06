begin;
CREATE TABLE if NOT EXISTS vehicles (
    id int NOT NULL PRIMARY KEY, 
    value jsonb NOT NULL
);
CREATE TABLE if NOT EXISTS groups (
    id int NOT NULL PRIMARY KEY,
    vehicle_id int,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    value jsonb NOT NULL
);
commit;
