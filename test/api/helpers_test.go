package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Llanos99/wf-engine/api"
	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/Llanos99/wf-engine/nodes"
	"github.com/Llanos99/wf-engine/script"
	"github.com/Llanos99/wf-engine/store"
)

// mockStore implements store.Store for testing
type mockStore struct {
	definitions *mockDefinitionRepo
	instances   *mockInstanceRepo
	approvals   *mockApprovalRepo
}

func newMockStore() *mockStore {
	return &mockStore{
		definitions: &mockDefinitionRepo{data: make(map[string]*domain.WorkflowDefinition)},
		instances:   &mockInstanceRepo{data: make(map[string]*domain.WorkflowInstance)},
		approvals:   &mockApprovalRepo{data: make(map[string]*domain.ApprovalRequest)},
	}
}

func (m *mockStore) Definitions() store.DefinitionRepository { return m.definitions }
func (m *mockStore) Instances() store.InstanceRepository     { return m.instances }
func (m *mockStore) Approvals() store.ApprovalRepository     { return m.approvals }
func (m *mockStore) Close() error                            { return nil }

// testServer creates a test server with mocked dependencies
func testServer(t *testing.T) (*api.Server, *mockStore) {
	t.Helper()
	mockStore := newMockStore()
	registry := nodes.NewRegistryWithBuiltin()
	scriptRunner := script.NewAdapter()
	server := api.NewServer(mockStore, registry, scriptRunner)
	return server, mockStore
}

// doRequest performs an HTTP request against the test server
func doRequest(t *testing.T, server *api.Server, method, path string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal body: %v", err)
		}
		reqBody = bytes.NewBuffer(data)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	server.Router().ServeHTTP(rec, req)
	return rec
}

// parseJSON parses the response body as JSON
func parseJSON(t *testing.T, rec *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(v); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
}

// assertStatus checks the response status code
func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if rec.Code != expected {
		t.Errorf("expected status %d, got %d: %s", expected, rec.Code, rec.Body.String())
	}
}

// mustJSON converts to json.RawMessage
func mustJSON(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}

// testDefinition creates a test workflow definition
func testDefinition() *domain.WorkflowDefinition {
	return &domain.WorkflowDefinition{
		ID:          "test-workflow",
		Version:     1,
		Name:        "Test Workflow",
		Description: "A test workflow",
		Nodes: map[string]domain.NodeDefinition{
			"start": {ID: "start", Type: domain.NodeTypeStart, Spec: mustJSON(map[string]string{"next": "end"})},
			"end":   {ID: "end", Type: domain.NodeTypeEnd, Spec: mustJSON(map[string]string{})},
		},
		InitialVariables: domain.Variables{"test": "value"},
	}
}

// mockDefinitionRepo is a mock implementation
type mockDefinitionRepo struct {
	data map[string]*domain.WorkflowDefinition
}

func (m *mockDefinitionRepo) Save(ctx context.Context, def *domain.WorkflowDefinition) error {
	key := def.ID + "-" + string(rune(def.Version))
	m.data[key] = def
	m.data[def.ID+"-latest"] = def
	return nil
}

func (m *mockDefinitionRepo) Get(ctx context.Context, id string, version int) (*domain.WorkflowDefinition, error) {
	key := id + "-" + string(rune(version))
	if def, ok := m.data[key]; ok {
		return def, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockDefinitionRepo) GetLatest(ctx context.Context, id string) (*domain.WorkflowDefinition, error) {
	key := id + "-latest"
	if def, ok := m.data[key]; ok {
		return def, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockDefinitionRepo) List(ctx context.Context) ([]*domain.WorkflowDefinition, error) {
	var result []*domain.WorkflowDefinition
	seen := make(map[string]bool)
	for _, def := range m.data {
		if !seen[def.ID] {
			result = append(result, def)
			seen[def.ID] = true
		}
	}
	return result, nil
}

func (m *mockDefinitionRepo) Delete(ctx context.Context, id string, version int) error {
	key := id + "-" + string(rune(version))
	delete(m.data, key)
	return nil
}

// mockInstanceRepo is a mock implementation
type mockInstanceRepo struct {
	data map[string]*domain.WorkflowInstance
	logs map[string][]domain.LogEntry
}

func (m *mockInstanceRepo) Create(ctx context.Context, inst *domain.WorkflowInstance) error {
	m.data[inst.ID] = inst
	return nil
}

func (m *mockInstanceRepo) Update(ctx context.Context, inst *domain.WorkflowInstance) error {
	m.data[inst.ID] = inst
	return nil
}

func (m *mockInstanceRepo) Get(ctx context.Context, id string) (*domain.WorkflowInstance, error) {
	if inst, ok := m.data[id]; ok {
		return inst, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockInstanceRepo) List(ctx context.Context, filter store.InstanceFilter) ([]*domain.WorkflowInstance, error) {
	var result []*domain.WorkflowInstance
	for _, inst := range m.data {
		if filter.DefinitionID != "" && inst.DefinitionID != filter.DefinitionID {
			continue
		}
		if filter.Status != "" && inst.Status != filter.Status {
			continue
		}
		result = append(result, inst)
	}
	return result, nil
}

func (m *mockInstanceRepo) Delete(ctx context.Context, id string) error {
	delete(m.data, id)
	return nil
}

func (m *mockInstanceRepo) SaveLog(ctx context.Context, instanceID string, log domain.LogEntry) error {
	if m.logs == nil {
		m.logs = make(map[string][]domain.LogEntry)
	}
	m.logs[instanceID] = append(m.logs[instanceID], log)
	return nil
}

func (m *mockInstanceRepo) GetLogs(ctx context.Context, instanceID string) ([]domain.LogEntry, error) {
	return m.logs[instanceID], nil
}

// mockApprovalRepo is a mock implementation
type mockApprovalRepo struct {
	data map[string]*domain.ApprovalRequest
}

func (m *mockApprovalRepo) Create(ctx context.Context, approval *domain.ApprovalRequest) error {
	m.data[approval.ID] = approval
	return nil
}

func (m *mockApprovalRepo) Get(ctx context.Context, id string) (*domain.ApprovalRequest, error) {
	if a, ok := m.data[id]; ok {
		return a, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockApprovalRepo) GetByInstance(ctx context.Context, instanceID string) ([]*domain.ApprovalRequest, error) {
	var result []*domain.ApprovalRequest
	for _, a := range m.data {
		if a.InstanceID == instanceID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockApprovalRepo) GetPendingByInstance(ctx context.Context, instanceID string) (*domain.ApprovalRequest, error) {
	for _, a := range m.data {
		if a.InstanceID == instanceID && a.Status == domain.ApprovalPending {
			return a, nil
		}
	}
	return nil, nil
}

func (m *mockApprovalRepo) Update(ctx context.Context, approval *domain.ApprovalRequest) error {
	m.data[approval.ID] = approval
	return nil
}

func (m *mockApprovalRepo) ListPending(ctx context.Context) ([]*domain.ApprovalRequest, error) {
	var result []*domain.ApprovalRequest
	for _, a := range m.data {
		if a.Status == domain.ApprovalPending {
			result = append(result, a)
		}
	}
	return result, nil
}

func (m *mockApprovalRepo) Delete(ctx context.Context, id string) error {
	delete(m.data, id)
	return nil
}
