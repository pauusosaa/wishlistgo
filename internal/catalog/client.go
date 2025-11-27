package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/httpx"
	"github.com/nmarsollier/commongo/log"
)

// Client permite comunicarse con el servicio de catálogo
type Client interface {
	GetArticle(articleID, token string) (*Article, error)
}

type client struct {
	log        log.LogRusEntry
	httpClient httpx.HTTPClient
	baseURL    string
}

// NewClient crea un nuevo cliente para el servicio de catálogo
func NewClient(log log.LogRusEntry, httpClient httpx.HTTPClient, baseURL string) Client {
	return &client{
		log:        log,
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// Article representa un artículo del catálogo
// Nota: CatalogGo retorna el formato simplificado (name, description, image en el nivel raíz)
type Article struct {
	ID          string  `json:"_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
	Enabled     bool    `json:"enabled"`
}

// GetArticle obtiene información de un artículo del catálogo
func (c *client) GetArticle(articleID, token string) (*Article, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/articles/"+articleID, nil)
	if err != nil {
		c.log.Error(err)
		return nil, errs.Invalid
	}

	req.Header.Add("Authorization", "Bearer "+token)

	// Agregar correlation ID si existe
	if corrId, ok := c.log.Data()[log.LOG_FIELD_CORRELATION_ID].(string); ok {
		req.Header.Add(log.LOG_FIELD_CORRELATION_ID, corrId)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error(err)
		return nil, errs.Invalid
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errs.NotFound
	}

	if resp.StatusCode != http.StatusOK {
		// Leer el cuerpo del error para debugging
		bodyBytes := make([]byte, 1024)
		n, _ := resp.Body.Read(bodyBytes)
		c.log.Error("catalog service returned status", resp.StatusCode, "body:", string(bodyBytes[:n]))
		return nil, errs.Invalid
	}

	var article Article
	if err := json.NewDecoder(resp.Body).Decode(&article); err != nil {
		c.log.Error(err)
		return nil, errs.Invalid
	}

	return &article, nil
}

