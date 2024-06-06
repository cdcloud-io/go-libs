package config

type Config struct {
	App struct {
		Name      string `yaml:"name"`
		Version   string `yaml:"version"`
		CommitSha string `yaml:"commit_sha"`
		BuildID   string `yaml:"build_id"`
		BuildDate string `yaml:"build_date"`
		Env       string `yaml:"env"`
		Debug     bool   `yaml:"debug"`
	} `yaml:"app"`
	Server struct {
		Host           string `yaml:"host"`
		Port           string `yaml:"port"`
		HealthEndpoint string `yaml:"health_endpoint"`
		InfoEndpoint   string `yaml:"info_endpoint"`
	} `yaml:"server"`
	Azure struct {
		AzStorageAccountName string   `yaml:"az_storage_account_name"`
		AzStorageAccountKey  string   `yaml:"az_storage_account_key"`
		Queues               []string `yaml:"queues"`
	} `yaml:"azure"`
}
