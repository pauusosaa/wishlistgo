package wishlist

// Service define las operaciones de dominio sobre Wishlist
type Service interface {
	GetByUserID(userID string) (*Wishlist, error)
	AddItem(userID, articleID string, notes *string) (*Item, error)
	RemoveItem(userID, articleID string) (bool, error)
}

type service struct {
	repo Repository
}

// NewService crea un nuevo servicio de dominio de wishlist
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetByUserID(userID string) (*Wishlist, error) {
	return s.repo.FindByUserID(userID)
}

func (s *service) AddItem(userID, articleID string, notes *string) (*Item, error) {
	w, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Evitar duplicados
	for _, it := range w.Items {
		if it.ArticleID == articleID {
			return it, nil
		}
	}

	item := &Item{
		ID:         "item-1", // TODO: usar UUID
		WishlistID: w.ID,
		ArticleID:  articleID,
		AddedAt:    "2024-01-15T10:30:00Z",
		Notes:      notes,
	}

	w.Items = append(w.Items, item)

	if err := s.repo.Save(w); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *service) RemoveItem(userID, articleID string) (bool, error) {
	return s.repo.RemoveItem(userID, articleID)
}


