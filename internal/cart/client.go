package cart

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/httpx"
	"github.com/nmarsollier/commongo/log"
)

// Client permite comunicarse con el servicio de carrito
type Client interface {
	AddArticle(articleID string, quantity int, token string) error
}

type client struct {
	log        log.LogRusEntry
	httpClient httpx.HTTPClient
	baseURL    string
}

// NewClient crea un nuevo cliente para el servicio de carrito
func NewClient(log log.LogRusEntry, httpClient httpx.HTTPClient, baseURL string) Client {
	return &client{
		log:        log,
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// AddArticleRequest representa la petición para agregar un artículo al carrito
type AddArticleRequest struct {
	ArticleID string `json:"articleId"`
	Quantity  int    `json:"quantity"`
}

// AddArticle agrega un artículo al carrito
func (c *client) AddArticle(articleID string, quantity int, token string) error {
	body := AddArticleRequest{
		ArticleID: articleID,
		Quantity:  quantity,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		c.log.Error(err)
		return errs.Invalid
	}

	req, err := http.NewRequest("POST", c.baseURL+"/cart/article", bytes.NewBuffer(jsonBody))
	if err != nil {
		c.log.Error(err)
		return errs.Invalid
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	// Agregar correlation ID si existe
	if corrId, ok := c.log.Data()[log.LOG_FIELD_CORRELATION_ID].(string); ok {
		req.Header.Add(log.LOG_FIELD_CORRELATION_ID, corrId)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Error(err)
		return errs.Invalid
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errs.NotFound
	}

	if resp.StatusCode == http.StatusBadRequest {
		return errs.NewValidation().Add("article", "No disponible")
	}

	if resp.StatusCode != http.StatusOK {
		c.log.Error("cart service returned status", resp.StatusCode)
		return errs.Invalid
	}

	return nil
}

