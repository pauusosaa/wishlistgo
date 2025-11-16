package env

import (
	"cmp"
	"os"

	"github.com/nmarsollier/commongo/strs"
)

// Configuration properties para Wishlist
type Configuration struct {
	ServerName string `json:"serverName"`
	Port       int    `json:"port"`
	FluentURL  string `json:"fluentUrl"`
}

var config *Configuration

// Get obtiene las variables de entorno del sistema
func Get() *Configuration {
	if config == nil {
		config = load()
	}

	return config
}

// load carga la configuración con defaults
func load() *Configuration {
	return &Configuration{
		ServerName: cmp.Or(os.Getenv("SERVER_NAME"), "wishlistgo"),
		Port:       cmp.Or(strs.AtoiZero(os.Getenv("PORT")), 3005),
		FluentURL:  cmp.Or(os.Getenv("FLUENT_URL"), "localhost:24224"),
	}
}


