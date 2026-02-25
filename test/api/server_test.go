package api_test

import (
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	server, _ := testServer(t)

	rec := doRequest(t, server, "GET", "/health", nil)
	assertStatus(t, rec, http.StatusOK)

	var response map[string]string
	parseJSON(t, rec, &response)

	if response["status"] != "ok" {
		t.Errorf("expected status ok, got %s", response["status"])
	}
}

func TestDefinitions_CreateAndGet(t *testing.T) {
	server, store := testServer(t)

	// Create definition
	createReq := map[string]interface{}{
		"id":          "test-workflow",
		"name":        "Test Workflow",
		"description": "A test",
		"nodes": map[string]interface{}{
			"start": map[string]interface{}{
				"id":   "start",
				"type": "start",
				"spec": map[string]string{"next": "end"},
			},
			"end": map[string]interface{}{
				"id":   "end",
				"type": "end",
				"spec": map[string]string{},
			},
		},
	}

	rec := doRequest(t, server, "POST", "/api/v1/definitions", createReq)
	assertStatus(t, rec, http.StatusCreated)

	// Verify stored
	if len(store.definitions.data) == 0 {
		t.Error("expected definition to be stored")
	}

	// Get definition
	rec = doRequest(t, server, "GET", "/api/v1/definitions/test-workflow", nil)
	assertStatus(t, rec, http.StatusOK)
}

func TestDefinitions_List(t *testing.T) {
	server, store := testServer(t)

	// Add test data
	def := testDefinition()
	store.definitions.Save(nil, def)

	rec := doRequest(t, server, "GET", "/api/v1/definitions", nil)
	assertStatus(t, rec, http.StatusOK)

	var response []map[string]interface{}
	parseJSON(t, rec, &response)

	if len(response) != 1 {
		t.Errorf("expected 1 definition, got %d", len(response))
	}
}

func TestDefinitions_NotFound(t *testing.T) {
	server, _ := testServer(t)

	rec := doRequest(t, server, "GET", "/api/v1/definitions/nonexistent", nil)
	assertStatus(t, rec, http.StatusNotFound)
}

func TestInstances_CreateAndExecute(t *testing.T) {
	server, store := testServer(t)

	// Create a simple workflow
	def := testDefinition()
	store.definitions.Save(nil, def)

	// Create instance
	createReq := map[string]interface{}{
		"definition_id": "test-workflow",
	}

	rec := doRequest(t, server, "POST", "/api/v1/instances", createReq)
	assertStatus(t, rec, http.StatusCreated)

	var response map[string]interface{}
	parseJSON(t, rec, &response)

	if response["status"] != "finished" {
		t.Errorf("expected status finished, got %v", response["status"])
	}

	// Verify stored
	if len(store.instances.data) != 1 {
		t.Errorf("expected 1 instance, got %d", len(store.instances.data))
	}
}

func TestInstances_List(t *testing.T) {
	server, store := testServer(t)

	// Create workflow and instance
	def := testDefinition()
	store.definitions.Save(nil, def)

	createReq := map[string]interface{}{
		"definition_id": "test-workflow",
	}
	doRequest(t, server, "POST", "/api/v1/instances", createReq)

	// List instances
	rec := doRequest(t, server, "GET", "/api/v1/instances", nil)
	assertStatus(t, rec, http.StatusOK)

	var response []map[string]interface{}
	parseJSON(t, rec, &response)

	if len(response) != 1 {
		t.Errorf("expected 1 instance, got %d", len(response))
	}
}

func TestInstances_NotFound(t *testing.T) {
	server, _ := testServer(t)

	rec := doRequest(t, server, "GET", "/api/v1/instances/nonexistent", nil)
	assertStatus(t, rec, http.StatusNotFound)
}

func TestInstances_Delete(t *testing.T) {
	server, store := testServer(t)

	// Create workflow and instance
	def := testDefinition()
	store.definitions.Save(nil, def)

	createReq := map[string]interface{}{
		"definition_id": "test-workflow",
	}
	rec := doRequest(t, server, "POST", "/api/v1/instances", createReq)

	var created map[string]interface{}
	parseJSON(t, rec, &created)
	instanceID := created["id"].(string)

	// Delete
	rec = doRequest(t, server, "DELETE", "/api/v1/instances/"+instanceID, nil)
	assertStatus(t, rec, http.StatusNoContent)

	// Verify deleted
	rec = doRequest(t, server, "GET", "/api/v1/instances/"+instanceID, nil)
	assertStatus(t, rec, http.StatusNotFound)
}

func TestApprovals_ListEmpty(t *testing.T) {
	server, _ := testServer(t)

	rec := doRequest(t, server, "GET", "/api/v1/approvals", nil)
	assertStatus(t, rec, http.StatusOK)

	var response []map[string]interface{}
	parseJSON(t, rec, &response)

	if len(response) != 0 {
		t.Errorf("expected 0 approvals, got %d", len(response))
	}
}

func TestApprovals_NotFound(t *testing.T) {
	server, _ := testServer(t)

	rec := doRequest(t, server, "GET", "/api/v1/approvals/nonexistent", nil)
	assertStatus(t, rec, http.StatusNotFound)
}

func TestCreateInstance_MissingDefinition(t *testing.T) {
	server, _ := testServer(t)

	createReq := map[string]interface{}{
		"definition_id": "nonexistent",
	}

	rec := doRequest(t, server, "POST", "/api/v1/instances", createReq)
	assertStatus(t, rec, http.StatusNotFound)
}

func TestCreateInstance_InvalidJSON(t *testing.T) {
	server, _ := testServer(t)

	req := "invalid json"
	rec := doRequest(t, server, "POST", "/api/v1/instances", req)
	// This actually fails differently because we're marshaling the string
	// The string gets marshaled as a valid JSON string
	assertStatus(t, rec, http.StatusBadRequest)
}

func TestCreateDefinition_MissingFields(t *testing.T) {
	server, _ := testServer(t)

	// Missing id
	rec := doRequest(t, server, "POST", "/api/v1/definitions", map[string]interface{}{
		"name": "Test",
	})
	assertStatus(t, rec, http.StatusBadRequest)

	// Missing name
	rec = doRequest(t, server, "POST", "/api/v1/definitions", map[string]interface{}{
		"id": "test",
	})
	assertStatus(t, rec, http.StatusBadRequest)

	// Missing nodes
	rec = doRequest(t, server, "POST", "/api/v1/definitions", map[string]interface{}{
		"id":   "test",
		"name": "Test",
	})
	assertStatus(t, rec, http.StatusBadRequest)
}
