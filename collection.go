package scaffold

import (
	"context"
	"errors"
	"time"

	"github.com/alexsobiek/scaffold/http"
	"github.com/alexsobiek/scaffold/query"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AccessFn is a callback function which is called before a document is accessed from the database.
// This function can be used to check if the user has access to the document.
type AccessFn[T any] func(context.Context, primitive.ObjectID) error

// ReadFn is a callback function which is called when a document is read from the database.
// This function can be used to modify the data before it is returned to the caller.
type ReadFn[T any] func(context.Context, primitive.ObjectID, *T) (*T, error)

// WriteFn is a callback function which is called before a document is written to the database.
// This function can be used to modify the data before it is written to the database.
type WriteFn[T any] func(context.Context, primitive.ObjectID, *T) (*T, error)

type UpdateFn[T any] func(context.Context, primitive.ObjectID, *T, *bson.M) (*bson.M, error)

type DeleteFn[T any] func(context.Context, primitive.ObjectID) error

type Collection interface {
	Name() string
	Slug() string
	inject(*mongo.Collection, *gin.RouterGroup)
}

type CollectionOpts[T any] struct {
	Name       string
	Slug       string
	Defaults   []Document[T]
	Access     AccessFn[T]
	Read       ReadFn[T]
	Write      WriteFn[T]
	Update     UpdateFn[T]
	Delete     DeleteFn[T]
	Middleware []gin.HandlerFunc
	Routes     []gin.RouteInfo
}

type C[T any] struct {
	name       string
	slug       string
	defaults   []Document[T]
	mc         *mongo.Collection
	access     AccessFn[T]
	read       ReadFn[T]
	write      WriteFn[T]
	update     UpdateFn[T]
	delete     DeleteFn[T]
	middleware []gin.HandlerFunc
	routes     []gin.RouteInfo
}

func NewCollection[T any](opts CollectionOpts[T]) *C[T] {
	if opts.Access == nil {
		opts.Access = func(_ context.Context, id primitive.ObjectID) error {
			return nil
		}
	}

	if opts.Read == nil {
		opts.Read = func(_ context.Context, id primitive.ObjectID, data *T) (*T, error) {
			return data, nil
		}
	}

	if opts.Write == nil {
		opts.Write = func(_ context.Context, _ primitive.ObjectID, data *T) (*T, error) {
			return data, nil
		}
	}

	if opts.Update == nil {
		opts.Update = func(_ context.Context, _ primitive.ObjectID, _ *T, updates *bson.M) (*bson.M, error) {
			return updates, nil
		}
	}

	if opts.Delete == nil {
		opts.Delete = func(_ context.Context, _ primitive.ObjectID) error {
			return nil
		}
	}

	return &C[T]{
		name:       opts.Name,
		slug:       opts.Slug,
		defaults:   opts.Defaults,
		access:     opts.Access,
		read:       opts.Read,
		write:      opts.Write,
		update:     opts.Update,
		delete:     opts.Delete,
		middleware: opts.Middleware,
		routes:     opts.Routes,
	}
}

func (c C[T]) Name() string {
	return c.name
}

func (c C[T]) Slug() string {
	return c.slug
}

func (c *C[T]) inject(mc *mongo.Collection, rg *gin.RouterGroup) {
	c.mc = mc

	for i := range c.defaults {
		doc := c.defaults[i]

		_, err := c.FindById(Context, doc.ID)

		if err != nil {
			if err != mongo.ErrNoDocuments {
				panic(err)
			}
		} else {
			continue
		}

		now := primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))

		if doc.Created == 0 {
			doc.Created = now
		}

		if doc.LastUpdated == 0 {
			doc.LastUpdated = now
		}

		_, err = c.mc.InsertOne(Context, doc)

		if err != nil {
			panic(err)
		}
	}

	rg.Use(c.middleware...)

	for i := range c.routes {
		route := c.routes[i]
		rg.Match([]string{route.Method}, route.Path, route.HandlerFunc)
	}

	rg.POST("/", c.handlePost)
	rg.GET("/:id", c.handleGetById)
	rg.PATCH("/:id", c.handlePatch)
}

func (c *C[T]) Insert(ctx context.Context, data T) (*Document[T], error) {
	doc := createDocument(c, data)

	d, err := c.write(ctx, doc.ID, doc.Data)

	if err != nil {
		return nil, err
	}

	doc.Data = d

	_, err = c.mc.InsertOne(ctx, doc)

	if err != nil {
		return nil, err
	}

	// Call read for any additional data processing
	doc.Data, err = c.read(ctx, doc.ID, doc.Data)

	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (c *C[T]) Find(ctx context.Context, query query.Query) (*Document[T], error) {
	var doc *Document[T]

	err := c.mc.FindOne(ctx, query.Filter()).Decode(&doc)

	if err != nil {
		return nil, err
	}

	doc.collection = c

	d, err := c.read(ctx, doc.ID, doc.Data)

	if err != nil {
		return nil, err
	}

	doc.Data = d

	return doc, nil
}

func (c *C[T]) FindById(ctx context.Context, id primitive.ObjectID) (*Document[T], error) {
	if err := c.access(ctx, id); err != nil {
		return nil, err
	}

	return c.Find(ctx, query.ID(id))
}

func (c *C[T]) handleGetById(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		http.BadRequest(ctx, errors.New("invalid id"))
		return
	}

	doc, err := c.FindById(ctx, id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.NotFound(ctx, nil)
		} else {
			http.Error(ctx, err)
		}
		return
	}

	http.Ok(ctx, doc)
}

func (c *C[T]) handlePost(ctx *gin.Context) {
	// Ensure Content-Type is application/json
	if ctx.GetHeader("Content-Type") != "application/json" {
		http.BadRequest(ctx, nil)
		return
	}

	var data T

	if err := ctx.BindJSON(&data); err != nil {
		http.BadRequest(ctx, err)
		return
	}

	doc, err := c.Insert(ctx, data)

	if err != nil {
		http.Error(ctx, err)
		return
	}

	// Call read for any additional data processing
	doc.Data, err = c.read(ctx, doc.ID, doc.Data)

	if err != nil {
		http.Error(ctx, err)
		return
	}

	http.Created(ctx, doc)
}

func (c *C[T]) handlePatch(ctx *gin.Context) {
	// Ensure Content-Type is application/json
	if ctx.GetHeader("Content-Type") != "application/json" {
		http.BadRequest(ctx, nil)
		return
	}

	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		http.BadRequest(ctx, errors.New("invalid id"))
		return
	}

	doc, err := c.FindById(ctx, id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.NotFound(ctx, nil)
		} else {
			http.Error(ctx, err)
		}
		return
	}

	var updates bson.M

	if err := ctx.BindJSON(&updates); err != nil {
		http.BadRequest(ctx, err)
		return
	}

	updated, err := c.update(ctx, doc.ID, doc.Data, &updates)

	if err != nil {
		http.Error(ctx, err)
		return
	}

	err = doc.SetMany(ctx, *updated)

	if err != nil {
		http.Error(ctx, err)
		return
	}

	http.Ok(ctx, doc)
}
