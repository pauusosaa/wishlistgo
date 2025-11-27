package wishlist

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wishlist struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string             `bson:"user_id" json:"user_id"`
	Items        []*Item            `bson:"items" json:"items"`
	CreationDate string             `bson:"creation_date" json:"creation_date"`
	UpdateDate   string             `bson:"update_date" json:"update_date"`
	TotalItems   int                `bson:"-" json:"total_items"`
}

type Item struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WishlistID primitive.ObjectID `bson:"wishlist_id" json:"wishlist_id"`
	ArticleID  string             `bson:"article_id" json:"article_id"`
	AddedAt    string             `bson:"added_at" json:"added_at"`
	Notes      *string            `bson:"notes,omitempty" json:"notes,omitempty"`

	//Campos que buscamos en el catálogo
	ArticleName     *string `bson:"-" json:"article_name,omitempty"`
	ArticlePrice    *float32 `bson:"-" json:"article_price,omitempty"`
	ArticleImage    *string `bson:"-" json:"article_image,omitempty"`
	IsAvailable     *bool   `bson:"-" json:"is_available,omitempty"`
	IsDisabled      *bool   `bson:"-" json:"is_disabled,omitempty"`
	Stock            *int    `bson:"-" json:"stock,omitempty"`
	DisabledReason   *string `bson:"-" json:"disabled_reason,omitempty"`
}


