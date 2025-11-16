package di

import (
	"github.com/pauusosaa/wishlistgo/internal/services"
	"github.com/pauusosaa/wishlistgo/internal/wishlist"
)

// Injector define las dependencias compartidas
type Injector interface {
	WishlistService() services.Service
}

type Deps struct {
	currWishlistSvc services.Service
}

// NewInjector crea un nuevo inyector
func NewInjector() Injector {
	return &Deps{}
}

func (d *Deps) WishlistService() services.Service {
	if d.currWishlistSvc != nil {
		return d.currWishlistSvc
	}

	repo := wishlist.NewInMemoryRepository()
	wSvc := wishlist.NewService(repo)
	d.currWishlistSvc = services.NewService(wSvc)

	return d.currWishlistSvc
}


