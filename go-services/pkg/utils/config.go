package utils

import (
    "log"

    "github.com/spf13/viper"
)

func LoadConfig(path string) error {
    viper.SetConfigFile(path)
    err := viper.ReadInConfig()
    if err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // Config file not found; handle appropriately
            log.Printf("No configuration file found at %s", path)
            return nil
        } else {
            // Config file was found but another error was produced
            return err
        }
    }
    return nil
}

func GetConfig(key string) interface{} {
    return viper.Get(key)
}
