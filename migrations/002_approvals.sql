-- Approval Requests
-- Stores pending approvals for human-in-the-loop workflows

CREATE TABLE IF NOT EXISTS approval_requests (
    id VARCHAR(255) PRIMARY KEY,
    instance_id VARCHAR(255) NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    node_id VARCHAR(255) NOT NULL,

    -- Approval metadata
    title VARCHAR(500) NOT NULL,
    description TEXT,
    approvers TEXT[],                    -- List of allowed approvers (usernames/emails)

    -- Status tracking
    status VARCHAR(50) NOT NULL DEFAULT 'pending',  -- pending, approved, rejected, expired
    decided_by VARCHAR(255),             -- Who approved/rejected
    decision_notes TEXT,                 -- Optional notes from approver

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,  -- Optional expiration
    decided_at TIMESTAMP WITH TIME ZONE
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_approval_requests_instance ON approval_requests(instance_id);
CREATE INDEX IF NOT EXISTS idx_approval_requests_status ON approval_requests(status);
CREATE INDEX IF NOT EXISTS idx_approval_requests_expires ON approval_requests(expires_at) WHERE status = 'pending';

-- View for pending approvals
CREATE OR REPLACE VIEW pending_approvals AS
SELECT
    ar.id,
    ar.instance_id,
    ar.node_id,
    ar.title,
    ar.description,
    ar.approvers,
    ar.created_at,
    ar.expires_at,
    wi.definition_id,
    wd.name as workflow_name
FROM approval_requests ar
JOIN workflow_instances wi ON ar.instance_id = wi.id
JOIN workflow_definitions wd ON wi.definition_id = wd.id AND wi.definition_version = wd.version
WHERE ar.status = 'pending'
  AND (ar.expires_at IS NULL OR ar.expires_at > NOW());
