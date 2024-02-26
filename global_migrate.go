package migrate

import (
	"fmt"
	"runtime"

	"go.mongodb.org/mongo-driver/mongo"
)

var globalMigrate = map[string]*Migrate{}

func internalRegister(name string, up MigrationFunc, down MigrationFunc, skip int) error {
	_, file, _, _ := runtime.Caller(skip)
	version, description, err := extractVersionDescription(file)
	if err != nil {
		return err
	}
	if hasVersion(globalMigrate[name].migrations, version) {
		return fmt.Errorf("migration with version %v already registered", version)
	}
	globalMigrate[name].migrations = append(globalMigrate[name].migrations, Migration{
		Version:     version,
		Description: description,
		Up:          up,
		Down:        down,
	})
	return nil
}

// Register performs migration registration.
// Use case of this function:
//
// - Create a file called like "1_setup_indexes.go" ("<version>_<comment>.go").
//
// - Use the following template inside:
//
//	 package migrations
//
//	 import (
//		 "go.mongodb.org/mongo-driver/bson"
//		 "go.mongodb.org/mongo-driver/mongo"
//		 "go.mongodb.org/mongo-driver/mongo/options"
//		 "github.com/jhlia0/mongo-migrate"
//	 )
//
//	 func init() {
//		 Register(func(db *mongo.Database) error {
//		 	 opt := options.Index().SetName("my-index")
//		 	 keys := bson.D{{"my-key", 1}}
//		 	 model := mongo.IndexModel{Keys: keys, Options: opt}
//		 	 _, err := db.Collection("my-coll").Indexes().CreateOne(context.TODO(), model)
//		 	 if err != nil {
//		 		 return err
//		 	 }
//		 	 return nil
//		 }, func(db *mongo.Database) error {
//		 	 _, err := db.Collection("my-coll").Indexes().DropOne(context.TODO(), "my-index")
//		 	 if err != nil {
//		 		 return err
//		 	 }
//		 	 return nil
//		 })
//	 }
func Register(name string, up MigrationFunc, down MigrationFunc) error {
	return internalRegister(name, up, down, 2)
}

// MustRegister acts like Register but panics on errors.
func MustRegister(name string, up, down MigrationFunc) {
	if err := internalRegister(name, up, down, 2); err != nil {
		panic(err)
	}
}

// RegisteredMigrations returns all registered migrations.
func RegisteredMigrations(name string) []Migration {
	ret := make([]Migration, len(globalMigrate[name].migrations))
	copy(ret, globalMigrate[name].migrations)
	return ret
}

// SetDatabase sets database for global migrate.
func SetDatabase(name string, db *mongo.Database) {
	if m, ok := globalMigrate[name]; ok {
		m.db = db
	}
}

// SetMigrationsCollection changes default collection name for migrations history.
func SetMigrationsCollection(migrationName, collectionName string) {
	globalMigrate[migrationName].SetMigrationsCollection(collectionName)
}

// Version returns current database version.
func Version(name string) (uint64, string, error) {
	return globalMigrate[name].Version()
}

// Up performs "up" migration using registered migrations.
// Detailed description available in Migrate.Up().
func Up(name string, n int) error {
	return globalMigrate[name].Up(n)
}

// Down performs "down" migration using registered migrations.
// Detailed description available in Migrate.Down().
func Down(name string, n int) error {
	return globalMigrate[name].Down(n)
}
