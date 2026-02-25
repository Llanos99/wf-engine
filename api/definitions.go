package api

import (
	"encoding/json"
	"net/http"

	"github.com/Llanos99/wf-engine/internal/domain"
)

// CreateDefinitionRequest is the request body for creating a definition
type CreateDefinitionRequest struct {
	ID               string                           `json:"id"`
	Version          int                              `json:"version"`
	Name             string                           `json:"name"`
	Description      string                           `json:"description"`
	Nodes            map[string]domain.NodeDefinition `json:"nodes"`
	InitialVariables domain.Variables                 `json:"initial_variables"`
}

// DefinitionResponse is the response for a definition
type DefinitionResponse struct {
	ID          string `json:"id"`
	Version     int    `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NodeCount   int    `json:"node_count"`
}

// toResponse converts a definition to a response
func definitionToResponse(def *domain.WorkflowDefinition) DefinitionResponse {
	return DefinitionResponse{
		ID:          def.ID,
		Version:     def.Version,
		Name:        def.Name,
		Description: def.Description,
		NodeCount:   len(def.Nodes),
	}
}

// listDefinitions handles GET /api/v1/definitions
func (s *Server) listDefinitions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	definitions, err := s.store.Definitions().List(ctx)
	if err != nil {
		InternalError(w, err.Error())
		return
	}

	response := make([]DefinitionResponse, len(definitions))
	for i, def := range definitions {
		response[i] = definitionToResponse(def)
	}

	JSON(w, http.StatusOK, response)
}

// createDefinition handles POST /api/v1/definitions
func (s *Server) createDefinition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateDefinitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid JSON: "+err.Error())
		return
	}

	if req.ID == "" {
		BadRequest(w, "id is required")
		return
	}
	if req.Name == "" {
		BadRequest(w, "name is required")
		return
	}
	if len(req.Nodes) == 0 {
		BadRequest(w, "nodes are required")
		return
	}

	def := &domain.WorkflowDefinition{
		ID:               req.ID,
		Version:          req.Version,
		Name:             req.Name,
		Description:      req.Description,
		Nodes:            req.Nodes,
		InitialVariables: req.InitialVariables,
	}

	if def.Version == 0 {
		def.Version = 1
	}

	if err := s.store.Definitions().Save(ctx, def); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, definitionToResponse(def))
}

// getDefinition handles GET /api/v1/definitions/{id}
func (s *Server) getDefinition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := urlParam(r, "id")

	def, err := s.store.Definitions().GetLatest(ctx, id)
	if err != nil {
		NotFound(w, "definition not found: "+id)
		return
	}

	JSON(w, http.StatusOK, def)
}

// getDefinitionVersion handles GET /api/v1/definitions/{id}/versions/{version}
func (s *Server) getDefinitionVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := urlParam(r, "id")

	version, err := urlParamInt(r, "version")
	if err != nil {
		BadRequest(w, "invalid version")
		return
	}

	def, err := s.store.Definitions().Get(ctx, id, version)
	if err != nil {
		NotFound(w, "definition not found")
		return
	}

	JSON(w, http.StatusOK, def)
}

// deleteDefinition handles DELETE /api/v1/definitions/{id}/versions/{version}
func (s *Server) deleteDefinition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := urlParam(r, "id")

	version, err := urlParamInt(r, "version")
	if err != nil {
		BadRequest(w, "invalid version")
		return
	}

	if err := s.store.Definitions().Delete(ctx, id, version); err != nil {
		NotFound(w, "definition not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
