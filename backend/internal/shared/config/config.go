package config

type Config struct {
	Port string
	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName: getEnv("DB_NAME", "postgres"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := 
}