package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	Postgres PostgresConfig
	GRPC     GRPCConfig
	App      AppConfig
}

type PostgresConfig struct {
	Host        string
	Port        int
	DBName      string
	User        string
	Password    string
	SSLMode     string
	ConnTimeout int
}

type GRPCConfig struct {
	Host    string
	Port    int
	Timeout int
}

type AppConfig struct {
	CostEncoding int
	SecretKey    string
}

func (p PostgresConfig) ConnectionURL() (string, error) {
	host := p.Host
	port := p.Port
	if port < 1 && port > 65536 {
		return "", fmt.Errorf("PSQL_PORT invalid")
	}
	host = host + ":" + strconv.Itoa(p.Port)

	urlBuilder := &url.URL{
		Scheme: "postgres",
		Host:   host,
		Path:   p.DBName,
	}

	if p.User == "" || p.Password == "" {
		return "", fmt.Errorf("PSQL_USER or PSQL_PASSWORD invalid")
	}
	urlBuilder.User = url.UserPassword(p.User, p.Password)

	query := urlBuilder.Query()
	connTimeout := p.ConnTimeout
	if connTimeout < 1 {
		return "", fmt.Errorf("PSQL_CONN_TIMEOUT invalid")
	}
	query.Add("connect_timeout", strconv.Itoa(p.ConnTimeout))

	if p.SSLMode != "disable" && p.SSLMode != "enable" {
		return "", fmt.Errorf("PSQL_SSL_MODE invalid")
	}
	query.Add("sslmode", p.SSLMode)

	urlBuilder.RawQuery = query.Encode()

	return urlBuilder.String(), nil
}

func (g GRPCConfig) ConnectionURL() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}

func LoadConfig(filenames ...string) error {
	return godotenv.Load(filenames...)
}

func New() Config {
	return Config{Postgres: PostgresConfig{
		Host:        getEnvAsString("PSQL_HOST", "localhost"),
		Port:        getEnvAsInt("PSQL_PORT", 5432),
		DBName:      getEnvAsString("PSQL_DB_NAME", "postgres"),
		User:        getEnvAsString("PSQL_USER", "postgres"),
		Password:    getEnvAsString("PSQL_PASSWORD", "postgres"),
		SSLMode:     getEnvAsString("PSQL_SSL_MODE", "disable"),
		ConnTimeout: getEnvAsInt("PSQL_CONN_TIMEOUT", 60),
	},
		GRPC: GRPCConfig{
			Host:    getEnvAsString("GRPC_HOST", "localhost"),
			Port:    getEnvAsInt("GRPC_PORT", 9090),
			Timeout: getEnvAsInt("GRPC_TIMEOUT", 60),
		},
		App: AppConfig{
			CostEncoding: getEnvAsInt("APP_COST_ENCODING", 8),
			SecretKey:    getEnvAsString("APP_SECRET_KEY", "my secret key"),
		}}
}

func getEnvAsString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnvAsString(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}
