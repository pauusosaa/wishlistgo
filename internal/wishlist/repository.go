package wishlist

// Repository define la interfaz de persistencia de wishlist
type Repository interface {
	FindByUserID(userID string) (*Wishlist, error)
	Save(w *Wishlist) error
	RemoveItem(userID, articleID string) (bool, error)
}

// inMemoryRepository es una implementación simple en memoria
// para desarrollo inicial y pruebas.
type inMemoryRepository struct {
	store map[string]*Wishlist // key: user_id
}

// NewInMemoryRepository crea un repositorio en memoria
func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		store: map[string]*Wishlist{},
	}
}

func (r *inMemoryRepository) FindByUserID(userID string) (*Wishlist, error) {
	w, exists := r.store[userID]
	if !exists {
		return &Wishlist{
			ID:           "wishlist-" + userID,
			UserID:       userID,
			Items:        []*Item{},
			CreationDate: "2024-01-15T10:30:00Z",
			UpdateDate:   "2024-01-15T10:30:00Z",
			TotalItems:   0,
		}, nil
	}
	w.TotalItems = len(w.Items)
	return w, nil
}

func (r *inMemoryRepository) Save(w *Wishlist) error {
	if w.Items != nil {
		w.TotalItems = len(w.Items)
	}
	r.store[w.UserID] = w
	return nil
}

func (r *inMemoryRepository) RemoveItem(userID, articleID string) (bool, error) {
	w, exists := r.store[userID]
	if !exists {
		return false, nil
	}

	newItems := make([]*Item, 0, len(w.Items))
	found := false
	for _, it := range w.Items {
		if it.ArticleID == articleID {
			found = true
			continue
		}
		newItems = append(newItems, it)
	}

	if !found {
		return false, nil
	}

	w.Items = newItems
	w.TotalItems = len(newItems)
	r.store[userID] = w
	return true, nil
}


