// internal/types/types.go
package types

import "errors"

// AppConfig represents an application configuration
type AppConfig struct {
    AppID     string
    AppKey    string
    AppSecret string
}

var ErrAppNotFound = errors.New("app not found")