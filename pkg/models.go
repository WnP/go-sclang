package models

// Query represent go-sclang HTTP query
type Query struct {
	// Code string to pass to sclang
	Code string
	// if true sclang return value is returned
	Stdout bool
	// if true kill sclang and sc-synth server
	Kill bool
	// if true kill sc-synth and reload sclang (sending SIGUSR1)
	Reload bool
}
