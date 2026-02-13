-- WF-Engine Database Schema
-- PostgreSQL

-- Workflow Definitions
-- Almacena los blueprints de workflows (inmutables una vez creados)
CREATE TABLE IF NOT EXISTS workflow_definitions (
    id VARCHAR(255) NOT NULL,
    version INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    definition JSONB NOT NULL,  -- El YAML parseado como JSON
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    PRIMARY KEY (id, version)
);

-- Índices para búsquedas comunes
CREATE INDEX IF NOT EXISTS idx_workflow_definitions_id ON workflow_definitions(id);
CREATE INDEX IF NOT EXISTS idx_workflow_definitions_name ON workflow_definitions(name);

-- Workflow Instances
-- Almacena las ejecuciones de workflows
CREATE TABLE IF NOT EXISTS workflow_instances (
    id VARCHAR(255) PRIMARY KEY,
    definition_id VARCHAR(255) NOT NULL,
    definition_version INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    current_node VARCHAR(255),
    variables JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    finished_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_definition
        FOREIGN KEY (definition_id, definition_version)
        REFERENCES workflow_definitions(id, version)
);

-- Índices para queries comunes
CREATE INDEX IF NOT EXISTS idx_workflow_instances_definition ON workflow_instances(definition_id);
CREATE INDEX IF NOT EXISTS idx_workflow_instances_status ON workflow_instances(status);
CREATE INDEX IF NOT EXISTS idx_workflow_instances_created_at ON workflow_instances(created_at DESC);

-- Execution Logs
-- Almacena los logs de ejecución de cada instancia
CREATE TABLE IF NOT EXISTS execution_logs (
    id SERIAL PRIMARY KEY,
    instance_id VARCHAR(255) NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    level VARCHAR(20) NOT NULL,
    node_id VARCHAR(255),
    message TEXT NOT NULL
);

-- Índice para obtener logs por instancia
CREATE INDEX IF NOT EXISTS idx_execution_logs_instance ON execution_logs(instance_id, timestamp);

-- Función para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para workflow_instances
DROP TRIGGER IF EXISTS update_workflow_instances_updated_at ON workflow_instances;
CREATE TRIGGER update_workflow_instances_updated_at
    BEFORE UPDATE ON workflow_instances
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Vista útil para ver instancias con info de definición
CREATE OR REPLACE VIEW workflow_instances_view AS
SELECT
    wi.id,
    wi.definition_id,
    wi.definition_version,
    wd.name as workflow_name,
    wi.status,
    wi.current_node,
    wi.variables,
    wi.created_at,
    wi.updated_at,
    wi.finished_at,
    EXTRACT(EPOCH FROM (COALESCE(wi.finished_at, NOW()) - wi.created_at)) as duration_seconds
FROM workflow_instances wi
JOIN workflow_definitions wd ON wi.definition_id = wd.id AND wi.definition_version = wd.version;
