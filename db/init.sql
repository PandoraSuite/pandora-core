CREATE DATABASE pandora;

\c pandora;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS service (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    version VARCHAR(11) NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT service_status_check CHECK (status IN ('active', 'deactivated', 'deprecated')),
    CONSTRAINT service_name_version_unique UNIQUE (name, version)
);

CREATE TABLE IF NOT EXISTS client (
    id SERIAL PRIMARY KEY,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT client_type_check CHECK (type IN ('organization', 'developer')),
    CONSTRAINT client_name_unique UNIQUE (name),
    CONSTRAINT client_email_unique UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS project (
    id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT project_client_id_fk FOREIGN KEY (client_id)
        REFERENCES client(id) ON DELETE CASCADE,

    CONSTRAINT project_status_check CHECK (status IN ('in_production', 'in_development', 'deactivated')),
    CONSTRAINT project_name_client_id_unique UNIQUE (name, client_id)
);

CREATE TABLE IF NOT EXISTS environment (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT environment_project_id_fk FOREIGN KEY (project_id)
        REFERENCES project(id) ON DELETE CASCADE,

    CONSTRAINT environment_status_check CHECK (status IN ('active', 'deactivated')),
    CONSTRAINT environment_name_project_id_unique UNIQUE (name, project_id)
);

CREATE TABLE IF NOT EXISTS api_key (
    id SERIAL PRIMARY KEY,
    environment_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NULL,
    last_used TIMESTAMP WITH TIME ZONE NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), 

    CONSTRAINT api_key_environment_id_fk FOREIGN KEY (environment_id)
        REFERENCES environment(id) ON DELETE CASCADE,

    CONSTRAINT api_key_status_check CHECK (status IN ('active', 'deactivated')),
    CONSTRAINT api_key_key_unique UNIQUE (key)
);

CREATE TABLE IF NOT EXISTS project_service (
    project_id INTEGER NOT NULL,
    service_id INTEGER NOT NULL,
    max_request INTEGER NULL,
    reset_frequency TEXT NULL,
    next_reset TIMESTAMP WITH TIME ZONE NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(), 

    PRIMARY KEY (project_id, service_id),
    
    CONSTRAINT project_service_project_id_fk FOREIGN KEY (project_id)
        REFERENCES project(id) ON DELETE CASCADE,
    CONSTRAINT project_service_service_id_fk FOREIGN KEY (service_id)
        REFERENCES service(id) ON DELETE CASCADE,

    CONSTRAINT project_service_reset_frequency_check CHECK (reset_frequency IN ('daily', 'weekly', 'biweekly', 'monthly'))
);

CREATE TABLE IF NOT EXISTS environment_service (
    environment_id INTEGER NOT NULL,
    service_id INTEGER NOT NULL,
    max_request INTEGER NULL,
    available_request INTEGER NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY (environment_id, service_id),

    CONSTRAINT environment_service_environment_id_fk FOREIGN KEY (environment_id)
        REFERENCES environment(id) ON DELETE CASCADE,
    CONSTRAINT environment_service_service_id_fk FOREIGN KEY (service_id)
        REFERENCES service(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS request_log (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    environment_id INTEGER NOT NULL,
    service_id INTEGER NOT NULL,
    api_key TEXT NOT NULL,
    start_point UUID NOT NULL,
    request_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    execution_status TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT request_log_environment_service_fk FOREIGN KEY (environment_id, service_id) 
	    REFERENCES environment_service(environment_id, service_id) ON DELETE CASCADE,
    CONSTRAINT request_log_start_point_fk FOREIGN KEY (start_point)
        REFERENCES request_log(id) ON DELETE CASCADE,

    CONSTRAINT request_log_execution_status_check CHECK (execution_status IN ('success', 'failed', 'unauthorized', 'server error'))
);

CREATE TABLE IF NOT EXISTS reservation (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    environment_id INTEGER NOT NULL,
    service_id INTEGER NOT NULL,
    api_key TEXT NOT NULL,
    request_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT reservation_environment_service_fk FOREIGN KEY (environment_id, service_id) 
	    REFERENCES environment_service(environment_id, service_id) ON DELETE CASCADE
);


CREATE INDEX IF NOT EXISTS idx_key ON api_key (key);
CREATE INDEX IF NOT EXISTS idx_request_log_api_key ON request_log (api_key);
CREATE INDEX IF NOT EXISTS idx_reservation_api_key ON reservation (api_key);
CREATE INDEX IF NOT EXISTS idx_start_point ON request_log (start_point);


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
