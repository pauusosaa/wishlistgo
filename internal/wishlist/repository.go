package wishlist

// Repository define la interfaz para persistir wishlists
type Repository interface {
	FindByUserID(userID string) (*Wishlist, error)
	Save(w *Wishlist) error
	RemoveItem(userID, articleID string) (bool, error)
}
