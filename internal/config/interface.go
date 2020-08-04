package config

type Provider interface {
	GetCodeRoot() string
	GetQuitExclusions() []string
}
