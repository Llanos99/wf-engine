package models

type Step struct {
	ID      string
	Name    string
	Execute func(*Context) error
	NextID  string
}
