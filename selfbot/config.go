package selfbot

import (
	"encoding/json"
	"errors"
)

// Config represents what can be configured about a selfbotsession. Usually loaded
// from a JSON Config file.
type Config struct {
	Token             string `json:"token"`
	Prefix            string `json:"prefix"`
	DefaultAsciiFont  string `json:"default_ascii_font"`
}

func NewConfigDefault(token string) Config {
	return Config{
		Token:            token,
		Prefix:           ">",
	}
}

func NewConfigFromJson(jsonBytes []byte) (Config, error) {
	// use the default Config for anything left out in the json
	config := NewConfigDefault("")
	if err := json.Unmarshal(jsonBytes, &config); err != nil {
		return Config{}, err
	}
	if config.Token == "" {
		return Config{}, errors.New("no token was specified")
	}
	return config, nil
}

func NewConfigsFromJson(jsonBytes []byte) ([]Config, error) {
	var configs []Config
	if err := json.Unmarshal(jsonBytes, &configs); err != nil {
		return nil, err
	}

	for _, config := range configs {
		if config.Token == "" {
			return nil, errors.New("one or more configs did not contain a token")
		}
	}

	return configs, nil
}
