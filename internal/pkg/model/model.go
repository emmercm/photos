package model

// Model represents structs that can be loaded and saved from DB
type Model interface {
	load() error
	Save() error
	Delete() error
}
