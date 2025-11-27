package wishlist

import (
	"context"
	"errors"
	"time"

	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrID = errs.NewValidation().Add("id", "Invalid")

type mongoRepository struct {
	log        log.LogRusEntry
	collection db.Collection
}

func NewMongoRepository(log log.LogRusEntry, collection db.Collection) Repository {
	return &mongoRepository{
		log:        log,
		collection: collection,
	}
}

type dbUserIDFilter struct {
	UserID string `bson:"user_id"`
}

func (r *mongoRepository) FindByUserID(userID string) (*Wishlist, error) {
	wishlist := &Wishlist{}
	filter := dbUserIDFilter{UserID: userID}

	if err := r.collection.FindOne(context.Background(), filter, wishlist); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		r.log.Error(err)
		return nil, err
	}

	wishlist.TotalItems = len(wishlist.Items)
	return wishlist, nil
}

func (r *mongoRepository) Save(w *Wishlist) error {
	if w.Items != nil {
		w.TotalItems = len(w.Items)
	}

	if w.ID == primitive.NilObjectID {
		w.ID = primitive.NewObjectID()
		_, err := r.collection.InsertOne(context.Background(), w)
		if err != nil {
			r.log.Error(err)
			return err
		}
		return nil
	}

	filter := bson.M{"_id": w.ID}
	_, err := r.collection.ReplaceOne(context.Background(), filter, w)
	if err != nil {
		r.log.Error(err)
		return err
	}

	return nil
}

func (r *mongoRepository) RemoveItem(userID, articleID string) (bool, error) {
	w, err := r.FindByUserID(userID)
	if err != nil {
		return false, err
	}
	if w == nil {
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
	w.UpdateDate = time.Now().UTC().Format(time.RFC3339)
	if err := r.Save(w); err != nil {
		return false, err
	}

	return true, nil
}

