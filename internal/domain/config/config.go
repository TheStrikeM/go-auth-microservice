package config

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Username string `yaml:"username" env-default:"thestrikem"`
	Password string `yaml:"password" env-default:"root"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Db       string `yaml:"db" env-default:"strikedb"`
}
