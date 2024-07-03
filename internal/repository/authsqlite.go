package repository

import (
	"errors"

	"github.com/subbbbbaru/go_final_project/configs"
)

type AuthEnv struct {
}

func NewAuthFromEnv() *AuthEnv {
	return &AuthEnv{}
}

func (auth *AuthEnv) GetPassword() (string, error) {
	config := configs.New()

	if len(config.AuthConf.Password) > 0 {
		return config.AuthConf.Password, nil
	}
	return "", errors.New("password not found")
}
