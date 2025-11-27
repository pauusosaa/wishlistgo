package wishlist

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	GetByUserID(userID string) (*Wishlist, error)
	AddItem(userID, articleID string, notes *string) (*Item, error)
	RemoveItem(userID, articleID string) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetByUserID(userID string) (*Wishlist, error) {
	w, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if w == nil {
		return &Wishlist{
			ID:           primitive.NilObjectID,
			UserID:       userID,
			Items:        []*Item{},
			CreationDate: "",
			UpdateDate:   "",
			TotalItems:   0,
		}, nil
	}
	return w, nil
}

func (s *service) AddItem(userID, articleID string, notes *string) (*Item, error) {
	w, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if w == nil {
		w = &Wishlist{
			ID:           primitive.NilObjectID,
			UserID:       userID,
			Items:        []*Item{},
			CreationDate: time.Now().UTC().Format(time.RFC3339),
			UpdateDate:   time.Now().UTC().Format(time.RFC3339),
		}
	}

	for _, it := range w.Items {
		if it.ArticleID == articleID {
			return it, nil
		}
	}

	item := &Item{
		ID:         primitive.NewObjectID(),
		WishlistID: w.ID,
		ArticleID:  articleID,
		AddedAt:    time.Now().UTC().Format(time.RFC3339),
		Notes:      notes,
	}

	w.Items = append(w.Items, item)
	w.UpdateDate = time.Now().UTC().Format(time.RFC3339)

	if err := s.repo.Save(w); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *service) RemoveItem(userID, articleID string) (bool, error) {
	w, err := s.repo.FindByUserID(userID)
	if err != nil {
		return false, err
	}
	if w == nil {
		return false, nil
	}

	return s.repo.RemoveItem(userID, articleID)
}
