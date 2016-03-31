package main

import (
    "github.com/scalingdata/gcfg"
)

type Config struct {
        Gerrit struct {
            User     string
            Password string
            Host     string
            CI       string
            Connections int
        }
        Color ColorTheme
        Project struct {
            Alias []string
            // redefine CI user on project basis
            CI []string
        }
}

func ReadConfig(file string) (*Config, error) {
    var config Config
    err := gcfg.ReadFileInto(&config, file)
    if err != nil {
        return nil, err
    }

    return &config, nil
}

