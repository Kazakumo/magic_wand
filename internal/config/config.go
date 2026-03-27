package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Log    LogConfig    `mapstructure:"log"`
	Ollama OllamaConfig `mapstructure:"ollama"`
	Qdrant QdrantConfig `mapstructure:"qdrant"`
	Ingest IngestConfig `mapstructure:"ingest"`
	RAG    RAGConfig    `mapstructure:"rag"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // "json" or "text"
}

type OllamaConfig struct {
	BaseURL    string `mapstructure:"base_url"`
	EmbedModel string `mapstructure:"embed_model"`
	ChatModel  string `mapstructure:"chat_model"`
	TimeoutSec int    `mapstructure:"timeout_sec"`
}

type QdrantConfig struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	Collection string `mapstructure:"collection"`
}

type IngestConfig struct {
	Workers         int `mapstructure:"workers"`
	ChunkSize       int `mapstructure:"chunk_size"`
	ChunkOverlap    int `mapstructure:"chunk_overlap"`
	WatchDebounceMs int `mapstructure:"watch_debounce_ms"`
}

type RAGConfig struct {
	TopK           int     `mapstructure:"top_k"`
	AlphaCosine    float64 `mapstructure:"alpha_cosine"`
	BetaKeyword    float64 `mapstructure:"beta_keyword"`
	GammaAST       float64 `mapstructure:"gamma_ast"`
	ScoreThreshold float64 `mapstructure:"score_threshold"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME/.magic_wand")

	setDefaults(v)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("config load: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}
	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")

	v.SetDefault("ollama.base_url", "http://localhost:11434")
	v.SetDefault("ollama.embed_model", "nomic-embed-text")
	v.SetDefault("ollama.chat_model", "llama3.2")
	v.SetDefault("ollama.timeout_sec", 30)

	v.SetDefault("qdrant.host", "localhost")
	v.SetDefault("qdrant.port", 6334)
	v.SetDefault("qdrant.collection", "magic_wand")

	v.SetDefault("ingest.workers", 4)
	v.SetDefault("ingest.chunk_size", 512)
	v.SetDefault("ingest.chunk_overlap", 64)
	v.SetDefault("ingest.watch_debounce_ms", 500)

	v.SetDefault("rag.top_k", 20)
	v.SetDefault("rag.alpha_cosine", 0.6)
	v.SetDefault("rag.beta_keyword", 0.3)
	v.SetDefault("rag.gamma_ast", 0.1)
	v.SetDefault("rag.score_threshold", 0.5)
}
