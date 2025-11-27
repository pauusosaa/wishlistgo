package di

import (
	"github.com/nmarsollier/commongo/db"
	"github.com/nmarsollier/commongo/httpx"
	"github.com/nmarsollier/commongo/log"
	"github.com/nmarsollier/commongo/security"
	"github.com/pauusosaa/wishlistgo/internal/cart"
	"github.com/pauusosaa/wishlistgo/internal/catalog"
	"github.com/pauusosaa/wishlistgo/internal/env"
	"github.com/pauusosaa/wishlistgo/internal/services"
	"github.com/pauusosaa/wishlistgo/internal/wishlist"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

// Singletons
var database *mongo.Database
var wishlistCollection db.Collection
var httpClient httpx.HTTPClient

// Injector define las dependencias compartidas
type Injector interface {
	Logger() log.LogRusEntry
	Database() *mongo.Database
	HttpClient() httpx.HTTPClient
	SecurityRepository() security.SecurityRepository
	SecurityService() security.SecurityService
	CatalogClient() catalog.Client
	CartClient() cart.Client
	WishlistCollection() db.Collection
	WishlistRepository() wishlist.Repository
	WishlistService() wishlist.Service
	WishlistAppService() services.Service
}

type Deps struct {
	currLog            log.LogRusEntry
	currHttpClient     httpx.HTTPClient
	currDatabase       *mongo.Database
	currSecRepo        security.SecurityRepository
	currSecSvc         security.SecurityService
	currCatalogClient  catalog.Client
	currCartClient     cart.Client
	currWishlistColl   db.Collection
	currWishlistRepo   wishlist.Repository
	currWishlistSvc    wishlist.Service
	currWishlistAppSvc services.Service
}

// NewInjector crea un nuevo inyector
func NewInjector(logger log.LogRusEntry) Injector {
	return &Deps{
		currLog: logger,
	}
}

func (d *Deps) Logger() log.LogRusEntry {
	return d.currLog
}

func (d *Deps) HttpClient() httpx.HTTPClient {
	if d.currHttpClient != nil {
		return d.currHttpClient
	}

	if httpClient != nil {
		return httpClient
	}

	httpClient = httpx.Get()
	return httpClient
}

func (d *Deps) SecurityRepository() security.SecurityRepository {
	if d.currSecRepo != nil {
		return d.currSecRepo
	}
	d.currSecRepo = security.NewSecurityRepository(d.Logger(), d.HttpClient(), env.Get().SecurityServerURL)
	return d.currSecRepo
}

func (d *Deps) SecurityService() security.SecurityService {
	if d.currSecSvc != nil {
		return d.currSecSvc
	}
	d.currSecSvc = security.NewSecurityService(d.Logger(), d.SecurityRepository())
	return d.currSecSvc
}

func (d *Deps) CatalogClient() catalog.Client {
	if d.currCatalogClient != nil {
		return d.currCatalogClient
	}
	d.currCatalogClient = catalog.NewClient(d.Logger(), d.HttpClient(), env.Get().CatalogServerURL)
	return d.currCatalogClient
}

func (d *Deps) CartClient() cart.Client {
	if d.currCartClient != nil {
		return d.currCartClient
	}
	d.currCartClient = cart.NewClient(d.Logger(), d.HttpClient(), env.Get().CartServerURL)
	return d.currCartClient
}

func (d *Deps) Database() *mongo.Database {
	if d.currDatabase != nil {
		return d.currDatabase
	}

	if database != nil {
		return database
	}

	database, err := db.NewDatabase(env.Get().MongoURL, "wishlist")
	if err != nil {
		d.currLog.Fatal(err)
		return nil
	}

	return database
}

func (d *Deps) WishlistCollection() db.Collection {
	if d.currWishlistColl != nil {
		return d.currWishlistColl
	}

	if wishlistCollection != nil {
		return wishlistCollection
	}

	wishlistCollection, err := db.NewCollection(d.Logger(), d.Database(), "wishlist", IsDbTimeoutError)
	if err != nil {
		d.currLog.Fatal(err)
		return nil
	}
	return wishlistCollection
}

func (d *Deps) WishlistRepository() wishlist.Repository {
	if d.currWishlistRepo != nil {
		return d.currWishlistRepo
	}
	d.currWishlistRepo = wishlist.NewMongoRepository(d.Logger(), d.WishlistCollection())
	return d.currWishlistRepo
}

func (d *Deps) WishlistService() wishlist.Service {
	if d.currWishlistSvc != nil {
		return d.currWishlistSvc
	}
	d.currWishlistSvc = wishlist.NewService(d.WishlistRepository())
	return d.currWishlistSvc
}

func (d *Deps) WishlistAppService() services.Service {
	if d.currWishlistAppSvc != nil {
		return d.currWishlistAppSvc
	}

	d.currWishlistAppSvc = services.NewService(d.Logger(), d.WishlistService(), d.CatalogClient(), d.CartClient())
	return d.currWishlistAppSvc
}

// IsDbTimeoutError función a llamar cuando se produce un error de db
func IsDbTimeoutError(err error) {
	if err == topology.ErrServerSelectionTimeout {
		database = nil
		wishlistCollection = nil
	}
}
