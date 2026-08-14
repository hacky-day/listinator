package database

import (
	"path/filepath"
	"testing"
)

// expectedSlugs are the type slugs the migrations seed. The frontend keeps a
// translation for each of them, see frontend/src/i18n/locales.
var expectedSlugs = []string{
	"misc",
	"fruit-vegetables-nuts",
	"canned-goods",
	"sauces-spices-dressings",
	"drinks-alcohol",
	"bakery",
	"spreads",
	"coffee-tea",
	"cereals-muesli",
	"pasta-rice",
	"cooking-baking",
	"meat-fish",
	"frozen",
	"dairy-chilled",
	"sweets-snacks",
	"household-baby-pets",
	"ready-meals-broth-sauce",
}

// initTestDB runs the migrations against a fresh database in a temporary
// directory. The admin bootstrap in Init is skipped by leaving the password
// empty.
func initTestDB(t *testing.T, dsn string) []Type {
	t.Helper()
	t.Setenv("LISTINATOR_ADMIN_PASSWORD", "")

	db, err := Init(dsn)
	if err != nil {
		t.Fatalf("unable to initialize database, %v", err)
	}

	ts := []Type{}
	if err := db.Find(&ts).Error; err != nil {
		t.Fatalf("unable to read types, %v", err)
	}
	return ts
}

func TestMigrateSeedsTypesWithSlug(t *testing.T) {
	ts := initTestDB(t, filepath.Join(t.TempDir(), "test.db"))

	if len(ts) != len(expectedSlugs) {
		t.Fatalf("got %v types, want %v", len(ts), len(expectedSlugs))
	}

	slugs := map[string]int{}
	for _, ty := range ts {
		if ty.Slug == "" {
			t.Errorf("type %q has no slug", ty.Name)
			continue
		}
		slugs[ty.Slug]++
		if slugs[ty.Slug] > 1 {
			t.Errorf("slug %q is used by more than one type", ty.Slug)
		}
	}

	for _, slug := range expectedSlugs {
		if slugs[slug] == 0 {
			t.Errorf("slug %q is missing", slug)
		}
	}
}

func TestMigrateIsIdempotent(t *testing.T) {
	dsn := filepath.Join(t.TempDir(), "test.db")

	initTestDB(t, dsn)
	ts := initTestDB(t, dsn)

	if len(ts) != len(expectedSlugs) {
		t.Fatalf("got %v types after migrating twice, want %v", len(ts), len(expectedSlugs))
	}
}
