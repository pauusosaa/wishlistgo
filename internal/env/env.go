package env

import (
	"cmp"
	"os"

	"github.com/nmarsollier/commongo/strs"
)

// Configuration properties para Wishlist
type Configuration struct {
	ServerName        string `json:"serverName"`
	Port              int    `json:"port"`
	MongoURL          string `json:"mongoUrl"`
	SecurityServerURL string `json:"securityServerUrl"`
	CatalogServerURL  string `json:"catalogUrl"`
	CartServerURL     string `json:"cartUrl"`
	FluentURL         string `json:"fluentUrl"`
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
		ServerName:        cmp.Or(os.Getenv("SERVER_NAME"), "wishlistgo"),
		Port:              cmp.Or(strs.AtoiZero(os.Getenv("PORT")), 3005),
		MongoURL:          cmp.Or(os.Getenv("MONGO_URL"), "mongodb://admin:admin123@localhost:27017/?authSource=admin"),
		SecurityServerURL: cmp.Or(os.Getenv("AUTH_SERVICE_URL"), "http://localhost:3000"),
		CatalogServerURL:  cmp.Or(os.Getenv("CATALOG_SERVICE_URL"), "http://localhost:3002"),
		CartServerURL:     cmp.Or(os.Getenv("CART_SERVICE_URL"), "http://localhost:3003"),
		FluentURL:         cmp.Or(os.Getenv("FLUENT_URL"), "localhost:24224"),
	}
}


