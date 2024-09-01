package utils

import (
    "net/http"
    "time"
)

var httpClient *http.Client

func GetHTTPClient() *http.Client {
    if httpClient == nil {
        httpClient = &http.Client{
            Timeout: time.Second * 10,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 100,
                IdleConnTimeout:     90 * time.Second,
            },
        }
    }
    return httpClient
}