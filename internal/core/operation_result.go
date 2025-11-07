package core

// OperationResult represents the result of a version manager operation
// Single Responsibility Principle: Only responsible for holding operation result data
type OperationResult struct {
	Success bool
	Message string
}
