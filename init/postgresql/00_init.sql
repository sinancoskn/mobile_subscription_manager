-- Create apps table first as it is referenced in other tables
CREATE TABLE apps (
    id SERIAL PRIMARY KEY,      -- Unique identifier for each app
    name VARCHAR(255) NOT NULL, -- Name of the app
    store VARCHAR(255) NOT NULL, -- Store (e.g., google, apple)
    created_at TIMESTAMPTZ DEFAULT NOW() -- Timestamp for creation
);

-- Create devices table and its partitions
CREATE TABLE devices (
    id BIGSERIAL,
    uid UUID NOT NULL,
    app_id INTEGER NOT NULL REFERENCES apps (id) ON DELETE CASCADE,
    language VARCHAR(10) NOT NULL,
    os SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (app_id, uid),
    UNIQUE (app_id, uid)
) PARTITION BY HASH (app_id);

CREATE TABLE devices_p0 PARTITION OF devices FOR VALUES WITH (MODULUS 2, REMAINDER 0);
CREATE TABLE devices_p1 PARTITION OF devices FOR VALUES WITH (MODULUS 2, REMAINDER 1);

-- Create subscriptions table and its partitions
CREATE TABLE subscriptions (
    id BIGSERIAL,                 
    uid UUID NOT NULL,            
    app_id INTEGER NOT NULL REFERENCES apps (id) ON DELETE CASCADE,      
    receipt VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,  
    expire_at TIMESTAMPTZ,        
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, app_id),
    UNIQUE (app_id, uid)     
) PARTITION BY HASH (app_id);

CREATE TABLE subscriptions_p0 PARTITION OF subscriptions FOR VALUES WITH (MODULUS 2, REMAINDER 0);
CREATE TABLE subscriptions_p1 PARTITION OF subscriptions FOR VALUES WITH (MODULUS 2, REMAINDER 1);

-- Create webhooks table
CREATE TABLE webhooks (
    id BIGSERIAL PRIMARY KEY,
    app_id INTEGER NOT NULL REFERENCES apps (id) ON DELETE CASCADE,
    url VARCHAR(1024) NOT NULL,
    default_headers JSONB,
    trigger_events JSONB NOT NULL,
    tried_count INT DEFAULT 0,
    last_attempt TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create manager_actions table
CREATE TABLE manager_actions (
    id BIGSERIAL PRIMARY KEY,
    expected_count BIGINT NOT NULL,
    will_be_processed_count BIGINT NOT NULL,
    max_batch INT NOT NULL,
    batch_count INT NOT NULL DEFAULT 0,
    completed_batch_count INT NOT NULL DEFAULT 0,
    triggered_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create batches table
CREATE TABLE batches (
    id BIGSERIAL PRIMARY KEY,
    action_id BIGINT NOT NULL REFERENCES manager_actions (id) ON DELETE CASCADE,
    start_index BIGINT NOT NULL,
    end_index BIGINT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    try_count INT NOT NULL DEFAULT 0,
    locked_by UUID DEFAULT NULL,
    locked_at TIMESTAMPTZ DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create workers table
CREATE TABLE workers (
    id BIGSERIAL PRIMARY KEY,
    worker_id UUID NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL,
    last_heartbeat TIMESTAMPTZ,
    action_id BIGINT DEFAULT NULL REFERENCES manager_actions (id) ON DELETE SET NULL,
    current_batch_id BIGINT DEFAULT NULL REFERENCES batches (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
