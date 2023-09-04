package config

import (
	"fmt"
	"os"
	"time"

	"github.com/pelletier/go-toml"
)

const DefaultFileName = "progress-bot-config.toml"

type (
	Config struct {
		Jira     *Jira     `toml:"jira"`
		Slack    *Slack    `toml:"slack"`
		OpenAI   *OpenAI   `toml:"openai"`
		Schedule *Schedule `toml:"schedule"`
	}

	Jira struct {
		User      string `toml:"user"`
		AccountID string `toml:"account_id"`
		BaseURL   string `toml:"base_url"`
		Token     string `toml:"token"`
		Project   string `toml:"project"`
	}

	Slack struct {
		StatusesEmoji []StatusEmoji `toml:"statuses_emoji"`
		Token         string        `toml:"token"`
		Channel       string        `toml:"channel"`
	}

	StatusEmoji struct {
		Name  string `toml:"name"`
		Emoji string `toml:"emoji"`
	}

	OpenAI struct {
		Token    string `toml:"token"`
		GPTModel string `toml:"gpt_model"`
	}

	Schedule struct {
		SinceDuration string `toml:"since_duration"`
	}
)

func LoadConfig(filepath string) (*Config, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %s. %s", filepath, err)
	}

	var cfg Config
	if err := toml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %s", err)
	}

	return &cfg, nil
}

func (cfg *Config) SinceDuration() (time.Duration, error) {
	d, err := time.ParseDuration(cfg.Schedule.SinceDuration)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %s. cfg duration: %s", err, cfg.Schedule.SinceDuration)
	}

	return d, nil
}
