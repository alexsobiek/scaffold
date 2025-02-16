# Scaffold
Scaffold is an ORM and HTTP library for building fast REST APIs. This project is barely in its alpha phase, expect many things to change.

## Example
```go
type SomeStruct struct {
	Name string `bson:"name" json:"name"`
}

func main() {
	c := scaffold.NewCollection(scaffold.CollectionOpts[SomeStruct]{
		Name: "Some Struct",
		Slug: "some-struct",
		Middleware: []gin.HandlerFunc{         		// define middleware
			func(c *gin.Context) {              	// Custom logic to lookup a user could go here
				c.Keys["user"] = "admin"        // set "user"
				c.Next()
			},
		},
        	// Read is called when we are reading from the database
		Read: func(ctx context.Context, id primitive.ObjectID, data SomeStruct) (SomeStruct, error) {
            		// look up the user
			fmt.Println(ctx.Value("user")) // admin
			return data, nil // return data without error
		},
        	// Write is called when we are writing a new object to the database
		Write: func(ctx context.Context, id primitive.ObjectID, data SomeStruct) (SomeStruct, error) {
			return data, http.ErrForbidden{} // Prevent writing, send 403 forbidden.
		},
	})

    	// Create new Scaffold instance
	s := scaffold.New(scaffold.ScaffoldOpts{
		Collections: []scaffold.Collection{c},  // Add collections
		MongoURI:    os.Getenv("MONGO_URI"),    // Set MongoDB URI
		Database:    "test",                    // Define database name
		Address:     os.Getenv("ADDRESS"),      // Set HTTP listen address
	})

    	// Run Scaffold
	if err := s.Run(context.Background()); err != nil {
		panic(err)
	}
}
```
### Create new SomeStruct entry
`curl -X POST -H "Content-Type: application/json" -d '{"name":"Test"}' http://localhost:3000/some-struct/`
Success:
```
{
  "data": {
    "id": "67b18bd1d31ddd889a15529d",
    "created": "2025-02-16T06:55:13.788Z",
    "last_updated": "2025-02-16T06:55:13.788Z",
    "document": {
      "name": "Test"
    }
  }
}
```
Forbidden:
```
{
  "error": "forbidden"
}
```
### Retrieve created entry
`curl http://localhost:3000/some-struct/67b18bd1d31ddd889a15529d`
Success:
```
{
  "data": {
    "id": "67b18bd1d31ddd889a15529d",
    "created": "2025-02-16T06:55:13.788Z",
    "last_updated": "2025-02-16T06:55:13.788Z",
    "document": {
      "name": "Test"
    }
  }
}
```
