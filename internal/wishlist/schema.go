package wishlist

// Wishlist representa la lista de deseos de un usuario
type Wishlist struct {
	ID           string        `bson:"_id,omitempty" json:"id"`
	UserID       string        `bson:"user_id" json:"user_id"`
	Items        []*Item       `bson:"items" json:"items"`
	CreationDate string        `bson:"creation_date" json:"creation_date"`
	UpdateDate   string        `bson:"update_date" json:"update_date"`
	TotalItems   int           `bson:"-" json:"total_items"`
}

// Item representa un artículo dentro de la wishlist
type Item struct {
	ID        string  `bson:"_id,omitempty" json:"id"`
	WishlistID string  `bson:"wishlist_id" json:"wishlist_id"`
	ArticleID string  `bson:"article_id" json:"article_id"`
	AddedAt   string  `bson:"added_at" json:"added_at"`
	Notes     *string `bson:"notes,omitempty" json:"notes,omitempty"`
}


