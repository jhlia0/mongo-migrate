package migrate

import (
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestGlobalMigrateSetGet(t *testing.T) {
	oldMigrate := globalMigrate
	defer func() {
		globalMigrate = oldMigrate
	}()
	db := &mongo.Database{}
	globalMigrate = map[string]*Migrate{}

	if globalMigrate["test"].db != db {
		t.Errorf("Unexpected non-equal dbs")
	}
	db2 := &mongo.Database{}
	SetDatabase("test", db2)
	if globalMigrate["test"].db != db2 {
		t.Errorf("Unexpected non-equal dbs")
	}
	SetMigrationsCollection("test", "test")
	if globalMigrate["test"].migrationsCollection != "test" {
		t.Errorf("Unexpected non-equal collections")
	}
}

func TestMigrationsRegistration(t *testing.T) {
	oldMigrate := globalMigrate
	defer func() {
		globalMigrate = oldMigrate
	}()
	globalMigrate = map[string]*Migrate{}

	err := Register("test", func(db *mongo.Database) error {
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
	if err != nil {
		t.Errorf("Unexpected register error: %v", err)
		return
	}
	registered := RegisteredMigrations("test")
	if len(registered) <= 0 || len(registered) > 1 {
		t.Errorf("Unexpected length of registered migrations")
		return
	}
	if registered[0].Version != 1 || registered[0].Description != "global_migrate_test" {
		t.Errorf("Unexpected version/description: %d %s", registered[0].Version, registered[0].Description)
	}

	err = Register("test", func(db *mongo.Database) error {
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
	if err == nil {
		t.Errorf("Unexpected nil error")
	}
}

func TestMigrationMustRegistration(t *testing.T) {
	oldMigrate := globalMigrate
	defer func() {
		globalMigrate = oldMigrate
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: %v", r)
		}
	}()
	globalMigrate = map[string]*Migrate{}
	MustRegister("test", func(db *mongo.Database) error {
		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
	registered := RegisteredMigrations("test")
	if len(registered) <= 0 || len(registered) > 1 {
		t.Errorf("Unexpected length of registered migrations")
		return
	}
	if registered[0].Version != 1 || registered[0].Description != "global_migrate_test" {
		t.Errorf("Unexpected version/description: %d %s", registered[0].Version, registered[0].Description)
	}
}
