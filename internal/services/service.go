package services

import (
	"github.com/pauusosaa/wishlistgo/internal/wishlist"
)

// Service expone los casos de uso de la Wishlist hacia REST / Rabbit / GraphQL
type Service interface {
	GetWishlist(userID string) (*wishlist.Wishlist, error)
	AddToWishlist(userID, articleID string, notes *string) (*wishlist.Item, error)
	RemoveFromWishlist(userID, articleID string) (bool, error)
	MoveToCart(userID, articleID string) error
}

type service struct {
	wishlist wishlist.Service
}

// NewService crea el servicio de aplicación para Wishlist
func NewService(w wishlist.Service) Service {
	return &service{
		wishlist: w,
	}
}

func (s *service) GetWishlist(userID string) (*wishlist.Wishlist, error) {
	return s.wishlist.GetByUserID(userID)
}

func (s *service) AddToWishlist(userID, articleID string, notes *string) (*wishlist.Item, error) {
	return s.wishlist.AddItem(userID, articleID, notes)
}

func (s *service) RemoveFromWishlist(userID, articleID string) (bool, error) {
	return s.wishlist.RemoveItem(userID, articleID)
}

// MoveToCart por ahora solo remueve de la wishlist.
// Más adelante integrará con CartGo via REST.
func (s *service) MoveToCart(userID, articleID string) error {
	_, err := s.wishlist.RemoveItem(userID, articleID)
	return err
}


