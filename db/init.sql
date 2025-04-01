CREATE DATABASE pandora;

\c pandora;

CREATE TABLE IF NOT EXISTS service (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    version VARCHAR(11) NOT NULL,
    status TEXT CHECK (status IN ('active', 'deactivated', 'deprecated')) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT unique_name_version UNIQUE(name, version)
);

CREATE TABLE IF NOT EXISTS client (
    id SERIAL PRIMARY KEY,
    type TEXT CHECK (type IN ('organization', 'developer')) NOT NULL,
    name TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS project (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    name TEXT NOT NULL,
    status TEXT CHECK (status IN ('in_production', 'in_development', 'deactivated')) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE,
    CONSTRAINT unique_name_client_id UNIQUE(name, client_id)
);

CREATE TABLE IF NOT EXISTS environment (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    name TEXT NOT NULL,
    status TEXT CHECK (status IN ('active', 'deactivated')) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE,
    CONSTRAINT unique_name_project_id UNIQUE(name, project_id)
);

CREATE TABLE IF NOT EXISTS api_key (
    id SERIAL PRIMARY KEY,
    environment_id INT NOT NULL,
    key TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NULL,
    last_used TIMESTAMP WITH TIME ZONE NULL,
    status TEXT CHECK (status IN ('active', 'deactivated')) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), 
    FOREIGN KEY (environment_id) REFERENCES environment(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_service (
    project_id INT NOT NULL,
    service_id INT NOT NULL,
    max_request INT NULL,
    reset_frequency TEXT CHECK (reset_frequency IN ('daily', 'weekly', 'biweekly', 'monthly')) NULL,
    next_reset TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), 
    PRIMARY KEY (project_id, service_id),
    FOREIGN KEY (project_id) REFERENCES project(id) ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS environment_service (
    environment_id INT NOT NULL,
    service_id INT NOT NULL,
    max_request INT NULL,
    available_request INT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (environment_id, service_id),
    FOREIGN KEY (environment_id) REFERENCES environment(id) ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES service(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS request_log (
    id TEXT PRIMARY KEY,
    environment_id INT NOT NULL,
    service_id INT NOT NULL,
    api_key TEXT NOT NULL,
    start_point TEXT NOT NULL,
    request_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    execution_status TEXT CHECK (execution_status IN ('success', 'failed', 'unauthorized', 'server error')) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    FOREIGN KEY (environment_id) REFERENCES environment_service(environment_id) ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES environment_service(service_id) ON DELETE CASCADE,
    FOREIGN KEY (start_point) REFERENCES request_log(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS reservations (
    id TEXT PRIMARY KEY,
    environment_id INT NOT NULL,
    service_id INT NOT NULL,
    api_key TEXT NOT NULL,
    environment TEXT NOT NULL,
    service TEXT NOT NULL,
    version varchar(11) NOT NULL,
    request_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (environment_id) REFERENCES environment_service(environment_id) ON DELETE CASCADE,
    FOREIGN KEY (service_id) REFERENCES environment_service(service_id) ON DELETE CASCADE,
);


CREATE INDEX IF NOT EXISTS idx_api_key_key ON api_key (key);
CREATE INDEX IF NOT EXISTS idx_request_log_api_key ON request_log (api_key);
CREATE INDEX IF NOT EXISTS idx_request_log_start_point ON request_log (start_point);


CREATE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER on_update_set_updated_at
    BEFORE UPDATE ON service
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER on_update_set_updated_at
    BEFORE UPDATE ON client
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER on_update_set_updated_at
    BEFORE UPDATE ON project
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER on_update_set_updated_at
    BEFORE UPDATE ON environment
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER on_update_set_updated_at
    BEFORE UPDATE ON api_key
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER on_update_set_updated_at
    BEFORE UPDATE ON project_service
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER on_update_set_updated_at
    BEFORE UPDATE ON environment_service
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();