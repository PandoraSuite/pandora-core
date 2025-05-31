CREATE DATABASE pandora;

\c pandora;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS service(
    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL,
    version VARCHAR(16) NOT NULL,
    CONSTRAINT service_name_version_unique UNIQUE (name, version),

    status TEXT NOT NULL,
    CONSTRAINT service_status_check
        CHECK (status IN ('enabled', 'disabled', 'deprecated')),

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS client(
    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL,
    CONSTRAINT client_name_unique UNIQUE (name),

    email TEXT NOT NULL,
    CONSTRAINT client_email_unique UNIQUE (email),

    type TEXT NOT NULL,
    CONSTRAINT client_type_check CHECK (type IN ('organization', 'developer')),

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project(
    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL,
    client_id INTEGER NOT NULL,
    CONSTRAINT project_client_id_fk
        FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE,
    CONSTRAINT project_name_client_id_unique UNIQUE (name, client_id),

    status TEXT NOT NULL,
    CONSTRAINT project_status_check
        CHECK (status IN ('enabled', 'disabled')),

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS environment(
    id SERIAL PRIMARY KEY,

    name TEXT NOT NULL,
    project_id INTEGER NOT NULL,
    CONSTRAINT environment_project_id_fk
        FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE,
    CONSTRAINT environment_name_project_id_unique UNIQUE (name, project_id),
    
    status TEXT NOT NULL,
    CONSTRAINT environment_status_check CHECK (status IN ('enabled', 'disabled')),

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS api_key(
    id SERIAL PRIMARY KEY,

    environment_id INTEGER NOT NULL,
    CONSTRAINT api_key_environment_id_fk
        FOREIGN KEY (environment_id) REFERENCES environment(id) ON DELETE CASCADE,

    key TEXT NOT NULL,
    CONSTRAINT api_key_key_unique UNIQUE (key),

    status TEXT NOT NULL,
    CONSTRAINT api_key_status_check
        CHECK (status IN ('enabled', 'disabled')),

    expires_at TIMESTAMPTZ,
    last_used TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project_service(
    project_id INTEGER NOT NULL,
    CONSTRAINT project_service_project_id_fk
        FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE,


    service_id INTEGER NOT NULL,
    CONSTRAINT project_service_service_id_fk
        FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE,

    reset_frequency TEXT,
    CONSTRAINT project_service_reset_frequency_check
        CHECK (reset_frequency IN ('daily', 'weekly', 'biweekly', 'monthly')),

    max_request INTEGER,
    next_reset TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS environment_service(
    environment_id INTEGER NOT NULL,
    CONSTRAINT environment_service_environment_id_fk
        FOREIGN KEY (environment_id) REFERENCES environment(id) ON DELETE CASCADE,

    service_id INTEGER NOT NULL,
    CONSTRAINT environment_service_service_id_fk
        FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE,

    CONSTRAINT environment_service_unique UNIQUE (environment_id, service_id),

    max_request INTEGER,
    available_request INTEGER,
    CONSTRAINT check_available_less_than_or_equal_max
        CHECK (available_request <= max_request),
    CONSTRAINT check_max_and_available_present_together
        CHECK (
            (max_request IS NOT NULL AND available_request IS NOT NULL)
            OR
            (max_request IS NULL AND available_request IS NULL)
        ),

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS request(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,

    -- Request chaining
    start_point UUID,
    CONSTRAINT request_start_point_fk
        FOREIGN KEY (start_point) REFERENCES request(id) ON DELETE CASCADE,

    -- API Key
    api_key TEXT NOT NULL,
    api_key_id INTEGER,
    CONSTRAINT request_api_key_id_fk
        FOREIGN KEY (api_key_id) REFERENCES api_key(id) ON DELETE SET NULL,

    -- Project
    project_name TEXT,
    project_id INTEGER,
    CONSTRAINT request_project_id_fk
        FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE SET NULL,

    -- Environment
    environment_name TEXT,
    environment_id INTEGER,
    CONSTRAINT request_environment_id_fk
        FOREIGN KEY (environment_id) REFERENCES environment(id) ON DELETE SET NULL,

    -- Service
    service_name TEXT NOT NULL,
    service_version VARCHAR(16) NOT NULL,
    service_id INTEGER,
    CONSTRAINT request_service_id_fk
        FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE,

    status_code INTEGER,
    execution_status TEXT NOT NULL,
    CONSTRAINT request_execution_status_check
        CHECK (
            execution_status IN (
                'success',
                'forwarded',
                'unauthorized',
                'client_error',
                'service_error'
            )
        ),
    CONSTRAINT request_status_code_required_check
        CHECK (
            execution_status NOT IN ('success', 'client_error', 'service_error')
            OR status_code IS NOT NULL
        ),

    unauthorized_reason TEXT,
    CONSTRAINT request_unauthorized_reason_check
        CHECK (
            unauthorized_reason IN (
                'API_KEY_INVALID',
                'QUOTA_EXCEEDED',
                'API_KEY_EXPIRED',
                'API_KEY_DISABLED',
                'SERVICE_MISMATCH',
                'SERVICE_DISABLED',
                'SERVICE_DEPRECATED',
                'SERVICE_NOT_ASSIGNED',
                'ENVIRONMENT_DISABLED'
            )
        ),

    request_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    ip_address TEXT NOT NULL,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    CONSTRAINT request_method_check
        CHECK (
            method IN (
                'GET',
                'HEAD',
                'POST',
                'PUT',
                'PATCH',
                'DELETE',
                'CONNECT',
                'OPTIONS',
                'TRACE'
            )
        ),
    
    metadata JSONB, 

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reservation(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,

    environment_id INTEGER NOT NULL,
    service_id INTEGER NOT NULL,
    CONSTRAINT reservation_environment_service_fk 
        FOREIGN KEY (environment_id, service_id) REFERENCES environment_service(environment_id, service_id) ON DELETE CASCADE,
    
    api_key TEXT NOT NULL,
    start_request_id UUID NOT NULL,
    request_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_key ON api_key(key);
CREATE INDEX IF NOT EXISTS idx_start_point ON request(start_point);
CREATE INDEX IF NOT EXISTS idx_request_api_key ON request(api_key);
CREATE INDEX IF NOT EXISTS idx_request_project_id ON request(project_id);
CREATE INDEX IF NOT EXISTS idx_request_service_id ON request(service_id);
CREATE INDEX IF NOT EXISTS idx_reservation_api_key ON reservation(api_key);
CREATE INDEX IF NOT EXISTS idx_request_environment_id ON request(environment_id);
