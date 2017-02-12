package config

type PlisConfig struct {
	Dependencies []PlisDependency `json:"dependencies"`
}
type PlisDependency struct {
	Repository string `json:"rep"`
	Branch string `json:"branch"`
}