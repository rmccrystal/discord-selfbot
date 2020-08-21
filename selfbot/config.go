package selfbot

import (
	"encoding/json"
	"errors"
)

// Config represents what can be configured about a selfbotsession. Usually loaded
// from a JSON config file.
type Config struct {
	Token            string `json:"token"`
	Prefix           string `json:"prefix"`
	DefaultAsciiFont string `json:"default_ascii_font"`
}

func NewConfigDefault(token string) Config {
	return Config{
		Token:  token,
		Prefix: ">",
		DefaultAsciiFont: "",
	}
}

func NewConfigFromJson(jsonBytes []byte) (Config, error) {
	// use the default config for anything left out in the json
	config := NewConfigDefault("")
	if err := json.Unmarshal(jsonBytes, &config); err != nil {
		return Config{}, err
	}
	if config.Token == "" {
		return Config{}, errors.New("no token in specified")
	}
	return config, nil
}
