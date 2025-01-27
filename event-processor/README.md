cmd/callback/main.go # entry point for callback project
cmd/worder/main.go # entry point for callback
internal/app/callback.go # app start logic for callback app. 
internal/models/ # model logic
internal/service/ # service logic


docker run --name redis -d -p 6379:6379 redis:latest

CREATE TABLE manager_actions (
    id BIGSERIAL PRIMARY KEY,               
    expected_count BIGINT NOT NULL,         
    will_be_processed_count BIGINT NOT NULL,
    max_batch INT NOT NULL,                
    batch_count INT NOT NULL DEFAULT 0,     -- Default value set to 0
    completed_batch_count INT NOT NULL DEFAULT 0, -- Default value set to 0
    triggered_at TIMESTAMPTZ NOT NULL,      
    status VARCHAR(20) DEFAULT 'pending',   -- Default value set to 'pending'
    created_at TIMESTAMPTZ DEFAULT NOW(),   -- Default value set to current timestamp
    updated_at TIMESTAMPTZ DEFAULT NOW()    -- Default value set to current timestamp
);

CREATE TABLE batches (
    id BIGSERIAL PRIMARY KEY,             
    action_id BIGINT NOT NULL,            
    start_index BIGINT NOT NULL,             
    end_index BIGINT NOT NULL,               
    status VARCHAR(20) DEFAULT 'pending', 
    try_count INT NOT NULL DEFAULT 0,
    locked_by UUID DEFAULT NULL,          
    locked_at TIMESTAMPTZ DEFAULT NULL,   
    created_at TIMESTAMPTZ DEFAULT NOW(), 
    updated_at TIMESTAMPTZ DEFAULT NOW(), 
    FOREIGN KEY (action_id) REFERENCES manager_actions (id) ON DELETE CASCADE
);

CREATE TABLE workers (
    id BIGSERIAL PRIMARY KEY,               -- Unique ID for the worker
    worker_id UUID NOT NULL UNIQUE,         -- UUID of the worker instance
    status VARCHAR(20) NOT NULL,            -- "idle", "processing", "stale"
    last_heartbeat TIMESTAMPTZ,             -- Last heartbeat from the worker
    action_id BIGINT DEFAULT NULL,          -- Reference to the current manager action
    current_batch_id BIGINT DEFAULT NULL,   -- Batch currently being processed
    created_at TIMESTAMPTZ DEFAULT NOW(),   -- Timestamp for creation
    updated_at TIMESTAMPTZ DEFAULT NOW(),   -- Timestamp for updates
    FOREIGN KEY (current_batch_id) REFERENCES batches (id) ON DELETE SET NULL
);
