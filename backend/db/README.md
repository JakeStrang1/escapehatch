# db package

This is where database integration code lives.

- This package should provide a common interface for interacting with the DB so that other packages don't need to import DB-specific packages like `mongo-driver` or `mgm`.
- This package should be fairly standalone. It shouldn't have dependencies on business-related packages like users or posts.

## Files

- `db.go` has helpers for reading/writing to the DB.
- `mgm.go` contains configuration code for the third-party `mgm` package, but it's low-level stuff.
- `registry.go` has some low level code used to swap "bson" struct tag for "db".
- `types.go` defines some helper types like `type db.M map[string]interface{}`.

## Errors

Functions in this package always return errors of type `internal/errors.Error`.


