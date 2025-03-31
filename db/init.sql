CREATE DATABASE pandora;

\c pandora;

CREATE TABLE IF NOT EXISTS service (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    version VARCHAR(11) NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT service_status_check CHECK (status IN ('active', 'deactivated', 'deprecated')),
    CONSTRAINT service_name_version_unique UNIQUE (name, version)
);

CREATE TABLE IF NOT EXISTS client (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    name TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT client_type_check CHECK (type IN ('organization', 'developer')),
    CONSTRAINT client_name_unique UNIQUE (name),
    CONSTRAINT client_email_unique UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS project (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT project_client_id_fk FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE,

    CONSTRAINT project_status_check CHECK (status IN ('in_production', 'in_development', 'deactivated')),
    CONSTRAINT project_name_client_id_unique UNIQUE (name, client_id)
);

CREATE TABLE IF NOT EXISTS environment (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT environment_project_id_fk FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE,

    CONSTRAINT environment_status_check CHECK (status IN ('active', 'deactivated')),
    CONSTRAINT environment_name_project_id_unique UNIQUE (name, project_id)
);

CREATE TABLE IF NOT EXISTS api_key (
    id SERIAL PRIMARY KEY,
    environment_id INT NOT NULL,
    key TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NULL,
    last_used TIMESTAMP WITH TIME ZONE NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT api_key_environment_id_fk FOREIGN KEY (environment_id) REFERENCES environment(id) ON DELETE CASCADE,

    CONSTRAINT api_key_status_check CHECK (status IN ('active', 'deactivated')),
    CONSTRAINT api_key_key_unique UNIQUE (key)
);

CREATE TABLE IF NOT EXISTS project_service (
    project_id INT NOT NULL,
    service_id INT NOT NULL,
    max_request INT NULL,
    reset_frequency TEXT NULL,
    next_reset TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY (project_id, service_id),
    
    CONSTRAINT project_service_project_id_fk FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE,
    CONSTRAINT project_service_service_id_fk FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE,

    CONSTRAINT project_service_reset_frequency_check CHECK (reset_frequency IN ('daily', 'weekly', 'biweekly', 'monthly'))
);

CREATE TABLE IF NOT EXISTS environment_service (
    environment_id INT NOT NULL,
    service_id INT NOT NULL,
    max_request INT NULL,
    available_request INT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY (environment_id, service_id),

    CONSTRAINT environment_service_environment_id_fk FOREIGN KEY (environment_id) REFERENCES environment(id) ON DELETE CASCADE,
    CONSTRAINT environment_service_service_id_fk FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS request_log (
    id SERIAL PRIMARY KEY,
    environment_id INT NOT NULL,
    service_id INT NOT NULL,
    api_key TEXT NOT NULL,
    request_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    execution_status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT request_log_environment_id_fk FOREIGN KEY (environment_id) REFERENCES environment(id) ON DELETE CASCADE,
    CONSTRAINT request_log_service_id_fk FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE,

    CONSTRAINT request_log_execution_status_check CHECK (execution_status IN ('success', 'failed', 'unauthorized', 'server error'))
);

CREATE INDEX IF NOT EXISTS idx_key ON api_key (key);
CREATE INDEX IF NOT EXISTS idx_api_key ON request_log (api_key);
