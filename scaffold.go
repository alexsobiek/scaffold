package scaffold

import (
	"context"
	"log"
	"os"

	"github.com/alexsobiek/scaffold/http"
)

type ScaffoldOpts struct {
	Collections []Collection
	MongoURI    string
	Database    string
	Address     string
	Logger      *log.Logger
}

type Scaffold struct {
	opts ScaffoldOpts
	db   *Database
	http *http.HttpServer
}

func New(opts ScaffoldOpts) *Scaffold {
	if opts.Logger == nil {
		opts.Logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	if opts.Address == "" {
		opts.Address = ":3000"
	}

	s := &Scaffold{
		opts: opts,
	}

	return s
}

func (s *Scaffold) Run(ctx context.Context) error {
	db, err := NewDatabase(ctx, s.opts.MongoURI, s.opts.Database)

	if err != nil {
		return err
	}

	s.db = db

	s.http = http.Create(s.opts.Logger, s.opts.Address)

	for _, c := range s.opts.Collections {
		c.inject(db.Collection(c.Slug()), s.http.Router().Group(c.Slug()))
	}

	s.http.Run()

	return nil
}