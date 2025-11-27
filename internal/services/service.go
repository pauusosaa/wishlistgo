package services

import (
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/log"
	"github.com/pauusosaa/wishlistgo/internal/cart"
	"github.com/pauusosaa/wishlistgo/internal/catalog"
	"github.com/pauusosaa/wishlistgo/internal/wishlist"
)

type Service interface {
	GetWishlist(userID, token string) (*wishlist.Wishlist, error)
	AddToWishlist(userID, articleID string, notes *string, token string) (*wishlist.Item, error)
	RemoveFromWishlist(userID, articleID string) (bool, error)
	MoveToCart(userID, articleID, token string) error
}

type service struct {
	log          log.LogRusEntry
	wishlist     wishlist.Service
	catalogClient catalog.Client
	cartClient   cart.Client
}

func NewService(log log.LogRusEntry, w wishlist.Service, catalogClient catalog.Client, cartClient cart.Client) Service {
	return &service{
		log:          log,
		wishlist:     w,
		catalogClient: catalogClient,
		cartClient:   cartClient,
	}
}

func (s *service) GetWishlist(userID, token string) (*wishlist.Wishlist, error) {
	w, err := s.wishlist.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, item := range w.Items {
		article, err := s.catalogClient.GetArticle(item.ArticleID, token)
		if err != nil {
			disabled := true
			reason := "Artículo no disponible"
			if err == errs.NotFound {
				reason = "El producto ya no está disponible"
			}
			item.IsDisabled = &disabled
			item.IsAvailable = &disabled
			item.DisabledReason = &reason
			stock := 0
			item.Stock = &stock
			continue
		}

		item.ArticleName = &article.Name
		item.ArticlePrice = &article.Price
		item.ArticleImage = &article.Image
		item.Stock = &article.Stock

		isAvailable := article.Enabled && article.Stock > 0
		item.IsAvailable = &isAvailable
		isDisabled := !isAvailable
		item.IsDisabled = &isDisabled

		if isDisabled {
			reason := "Sin stock disponible"
			if !article.Enabled {
				reason = "Producto deshabilitado"
			}
			item.DisabledReason = &reason
		}
	}

	return w, nil
}

func (s *service) AddToWishlist(userID, articleID string, notes *string, token string) (*wishlist.Item, error) {
	_, err := s.catalogClient.GetArticle(articleID, token)
	if err != nil {
		if err == errs.NotFound {
			return nil, errs.NotFound
		}
		return nil, errs.Invalid
	}

	w, err := s.wishlist.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if w != nil {
		for _, it := range w.Items {
			if it.ArticleID == articleID {
				return nil, errs.NewValidation().Add("article_id", "El artículo ya está en la wishlist")
			}
		}
	}

	return s.wishlist.AddItem(userID, articleID, notes)
}

func (s *service) RemoveFromWishlist(userID, articleID string) (bool, error) {
	return s.wishlist.RemoveItem(userID, articleID)
}

func (s *service) MoveToCart(userID, articleID, token string) error {
	article, err := s.catalogClient.GetArticle(articleID, token)
	if err != nil {
		if err == errs.NotFound {
			return errs.NotFound
		}
		return errs.Invalid
	}

	if !article.Enabled || article.Stock <= 0 {
		return errs.NewValidation().Add("article", "El artículo no está disponible")
	}

	err = s.cartClient.AddArticle(articleID, 1, token)
	if err != nil {
		return err
	}

	_, err = s.wishlist.RemoveItem(userID, articleID)
	return err
}
