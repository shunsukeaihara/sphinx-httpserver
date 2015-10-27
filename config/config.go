package config

import (
	"errors"
	"io"
	"io/ioutil"

	ps "github.com/shunsukeaihara/go-pocketsphinx"
	"golang.org/x/net/context"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	PSConfig map[string]ps.Config `yaml:"ps_config"`
	Env      string
}

type key int

const configKey key = 0

type environments map[string]Config

func Load(r io.Reader, env string) (Config, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Config{}, err
	}

	var all environments
	err = yaml.Unmarshal(data, &all)
	if err != nil {
		return Config{}, err
	}

	cfg, ok := all[env]
	if !ok {
		return Config{}, errors.New("no such environment: " + env)
	}
	cfg.Env = env

	return cfg, nil
}

func NewContext(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, configKey, cfg)
}

func FromContext(ctx context.Context) (Config, bool) {
	cfg, ok := ctx.Value(configKey).(Config)
	return cfg, ok
}
