package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Vector   VectorConfig   `mapstructure:"vector"`
	Log      LogConfig      `mapstructure:"log"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Filename string `mapstructure:"filename"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"driver"`
	Source string `mapstructure:"source"`
}

type LLMConfig struct {
	Provider string            `mapstructure:"provider" json:"provider"`
	APIKey   string            `mapstructure:"api_key" json:"api_key"`
	BaseURL  string            `mapstructure:"base_url" json:"base_url"`
	Model          string            `mapstructure:"model" json:"model"`
	EmbeddingModel string            `mapstructure:"embedding_model" json:"embedding_model"`
	Keys           map[string]string `mapstructure:"keys" json:"keys"` // 支持多 Key 配置
}

type VectorConfig struct {
	Provider string `mapstructure:"provider"` // 如 milvus, pinecone, sqlite (local)
	Address  string `mapstructure:"address"`
	APIKey   string `mapstructure:"api_key"`
}

var GlobalConfig Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set aliases to handle both underscore and no-underscore keys
	viper.RegisterAlias("llm.apikey", "llm.api_key")
	viper.RegisterAlias("llm.baseurl", "llm.base_url")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.driver", "pgsql")
	viper.SetDefault("database.source", "host=localhost user=postgres password=postgres dbname=ai_novel port=5432 sslmode=disable")
	viper.SetDefault("llm.provider", "openai")
	viper.SetDefault("llm.api_key", "")
	viper.SetDefault("llm.base_url", "")
	viper.SetDefault("llm.model", "gpt-3.5-turbo")
	viper.SetDefault("llm.embedding_model", "text-embedding-ada-002")
	viper.SetDefault("vector.provider", "pgsql")
	viper.SetDefault("vector.address", "host=localhost user=postgres password=postgres dbname=ai_novel port=5432 sslmode=disable")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.filename", "app.log")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("Config file not found, using defaults and environment variables")
		} else {
			// Config file was found but another error produced
			return fmt.Errorf("fatal error config file: %w", err)
		}
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}

	return nil
}

func SaveConfig() error {
	// Sync GlobalConfig back to viper
	viper.Set("server.port", GlobalConfig.Server.Port)
	viper.Set("server.mode", GlobalConfig.Server.Mode)
	viper.Set("database.driver", GlobalConfig.Database.Driver)
	viper.Set("database.source", GlobalConfig.Database.Source)
	viper.Set("llm.provider", GlobalConfig.LLM.Provider)
	viper.Set("llm.api_key", GlobalConfig.LLM.APIKey)
	viper.Set("llm.base_url", GlobalConfig.LLM.BaseURL)
	viper.Set("llm.model", GlobalConfig.LLM.Model)
	viper.Set("llm.embedding_model", GlobalConfig.LLM.EmbeddingModel)
	viper.Set("llm.keys", GlobalConfig.LLM.Keys)
	viper.Set("vector.provider", GlobalConfig.Vector.Provider)
	viper.Set("vector.address", GlobalConfig.Vector.Address)
	viper.Set("vector.api_key", GlobalConfig.Vector.APIKey)
	viper.Set("log.level", GlobalConfig.Log.Level)
	viper.Set("log.filename", GlobalConfig.Log.Filename)

	// Try to write to config.yaml
	return viper.WriteConfig()
}
