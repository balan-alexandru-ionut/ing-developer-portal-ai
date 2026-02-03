package config

type Config struct {
	Vertex VertexAIConfig `yaml:"vertex"`
}

type VertexAIConfig struct {
	Project   ProjectConfig   `yaml:"project"`
	DataStore DataStoreConfig `yaml:"dataStore"`
	Model     AIModelConfig   `yaml:"model"`
}

type ProjectConfig struct {
	Id string `yaml:"id"`
}

type DataStoreConfig struct {
	Id       string `yaml:"id"`
	Location string `yaml:"location"`
}

type AIModelConfig struct {
	Name            string  `yaml:"name"`
	Location        string  `yaml:"location"`
	Temperature     float32 `yaml:"temperature"`
	MaxOutputTokens int32   `yaml:"maxOutputTokens"`
}
