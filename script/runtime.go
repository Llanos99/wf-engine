package script

import (
	"fmt"

	"github.com/Llanos99/wf-engine/internal/domain"
	"github.com/dop251/goja"
)

// Runtime ejecuta scripts JavaScript con acceso a las APIs del workflow
type Runtime struct {
	vm *goja.Runtime
}

// NewRuntime crea un nuevo runtime de JavaScript
func NewRuntime() *Runtime {
	return &Runtime{}
}

// Execute ejecuta un script con acceso al contexto del workflow
func (r *Runtime) Execute(script string, vars *domain.Variables, instance *domain.WorkflowInstance) (interface{}, error) {
	vm := goja.New()

	// Configurar APIs
	if err := r.setupAPIs(vm, vars, instance); err != nil {
		return nil, fmt.Errorf("failed to setup APIs: %w", err)
	}

	// Ejecutar el script
	result, err := vm.RunString(script)
	if err != nil {
		return nil, fmt.Errorf("script error: %w", err)
	}

	return result.Export(), nil
}

// ExecuteBool ejecuta un script y espera un boolean como resultado
func (r *Runtime) ExecuteBool(script string, vars *domain.Variables, instance *domain.WorkflowInstance) (bool, error) {
	result, err := r.Execute(script, vars, instance)
	if err != nil {
		return false, err
	}

	boolResult, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("script did not return a boolean, got: %T", result)
	}

	return boolResult, nil
}

// setupAPIs configura todas las APIs disponibles en el script
func (r *Runtime) setupAPIs(vm *goja.Runtime, vars *domain.Variables, instance *domain.WorkflowInstance) error {
	// API: workflow
	workflow := vm.NewObject()
	workflow.Set("vars", r.createVarsAPI(vm, vars))
	workflow.Set("instanceId", instance.ID)
	workflow.Set("definitionId", instance.DefinitionID)
	vm.Set("workflow", workflow)

	// API: log
	vm.Set("log", r.createLogAPI(vm, instance))

	// API: http
	vm.Set("http", r.createHTTPAPI(vm))

	// API: db (placeholder)
	vm.Set("db", r.createDBAPI(vm))

	// API: console (para debugging)
	console := vm.NewObject()
	console.Set("log", func(call goja.FunctionCall) goja.Value {
		args := make([]interface{}, len(call.Arguments))
		for i, arg := range call.Arguments {
			args[i] = arg.Export()
		}
		fmt.Println(args...)
		return goja.Undefined()
	})
	vm.Set("console", console)

	return nil
}

// createVarsAPI crea la API para acceder a las variables del workflow
func (r *Runtime) createVarsAPI(vm *goja.Runtime, vars *domain.Variables) *goja.Object {
	varsAPI := vm.NewObject()

	// get(key) - obtiene una variable
	varsAPI.Set("get", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}
		key := call.Arguments[0].String()
		val, ok := vars.Get(key)
		if !ok {
			return goja.Undefined()
		}
		return vm.ToValue(val)
	})

	// set(key, value) - establece una variable
	varsAPI.Set("set", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 2 {
			return goja.Undefined()
		}
		key := call.Arguments[0].String()
		value := call.Arguments[1].Export()
		(*vars)[key] = value
		return goja.Undefined()
	})

	// getAll() - obtiene todas las variables
	varsAPI.Set("getAll", func(call goja.FunctionCall) goja.Value {
		return vm.ToValue(map[string]interface{}(*vars))
	})

	// has(key) - verifica si existe una variable
	varsAPI.Set("has", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return vm.ToValue(false)
		}
		key := call.Arguments[0].String()
		_, ok := vars.Get(key)
		return vm.ToValue(ok)
	})

	// delete(key) - elimina una variable
	varsAPI.Set("delete", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}
		key := call.Arguments[0].String()
		delete(*vars, key)
		return goja.Undefined()
	})

	return varsAPI
}

// createLogAPI crea la API de logging
func (r *Runtime) createLogAPI(vm *goja.Runtime, instance *domain.WorkflowInstance) *goja.Object {
	logAPI := vm.NewObject()

	createLogFn := func(level string) func(goja.FunctionCall) goja.Value {
		return func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) < 1 {
				return goja.Undefined()
			}
			message := call.Arguments[0].String()
			instance.AddLog(level, instance.CurrentNode, message)
			fmt.Printf("[%s] %s: %s\n", level, instance.CurrentNode, message)
			return goja.Undefined()
		}
	}

	logAPI.Set("debug", createLogFn("debug"))
	logAPI.Set("info", createLogFn("info"))
	logAPI.Set("warn", createLogFn("warn"))
	logAPI.Set("error", createLogFn("error"))

	return logAPI
}

// createHTTPAPI crea la API para hacer llamadas HTTP
func (r *Runtime) createHTTPAPI(vm *goja.Runtime) *goja.Object {
	httpAPI := vm.NewObject()

	// Por ahora, placeholder que simula HTTP
	// TODO: Implementar llamadas HTTP reales
	httpAPI.Set("get", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}
		url := call.Arguments[0].String()
		fmt.Printf("[HTTP] GET %s\n", url)

		// Simular respuesta
		response := vm.NewObject()
		response.Set("status", 200)
		response.Set("body", map[string]interface{}{
			"message": "simulated response",
			"url":     url,
		})
		return response
	})

	httpAPI.Set("post", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}
		url := call.Arguments[0].String()
		var body interface{}
		if len(call.Arguments) > 1 {
			body = call.Arguments[1].Export()
		}
		fmt.Printf("[HTTP] POST %s - %v\n", url, body)

		response := vm.NewObject()
		response.Set("status", 200)
		response.Set("body", map[string]interface{}{
			"message": "simulated response",
			"url":     url,
		})
		return response
	})

	return httpAPI
}

// createDBAPI crea la API para acceder a la base de datos
func (r *Runtime) createDBAPI(vm *goja.Runtime) *goja.Object {
	dbAPI := vm.NewObject()

	// Placeholder para futuras operaciones de BD
	// TODO: Implementar conexión real a BD
	dbAPI.Set("query", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}
		query := call.Arguments[0].String()
		fmt.Printf("[DB] Query: %s\n", query)

		// Simular resultado vacío
		return vm.ToValue([]interface{}{})
	})

	dbAPI.Set("execute", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			return goja.Undefined()
		}
		query := call.Arguments[0].String()
		fmt.Printf("[DB] Execute: %s\n", query)

		result := vm.NewObject()
		result.Set("rowsAffected", 0)
		return result
	})

	return dbAPI
}
