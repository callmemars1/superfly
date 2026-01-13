-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS apps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(63) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    
    -- For Feature 1: Pre-built image deployment
    image TEXT NOT NULL,
    
    -- Runtime config
    port INTEGER NOT NULL DEFAULT 8080,
    replicas INTEGER NOT NULL DEFAULT 1,
    cpu_limit VARCHAR(10) NOT NULL DEFAULT '500m',
    memory_limit VARCHAR(10) NOT NULL DEFAULT '256Mi',
    
    -- Domain config
    domain VARCHAR(255) UNIQUE,
    
    -- Health check
    health_check_path VARCHAR(255) NOT NULL DEFAULT '/',
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    
    -- Metadata
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_deployed_at TIMESTAMPTZ
);

-- Index for faster lookups
CREATE INDEX idx_apps_slug ON apps(slug);
CREATE INDEX idx_apps_status ON apps(status);
CREATE INDEX idx_apps_domain ON apps(domain) WHERE domain IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS apps;
-- +goose StatementEnd
