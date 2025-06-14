package shroom

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bouncepaw/mycorrhiza/internal/cfg"
	"github.com/bouncepaw/mycorrhiza/internal/files"
	"github.com/bouncepaw/mycorrhiza/internal/hyphae"
)

func setupTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "myco")
	if err != nil {
		t.Fatal(err)
	}
	cfg.WikiDir = dir
	if err := files.PrepareWikiRoot(); err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestRenamingPairsIgnoresEmptyEntries(t *testing.T) {
	setupTestDir(t)

	emptyBase := hyphae.ByName("").(*hyphae.EmptyHypha)
	empty := hyphae.ExtendEmptyToTextual(emptyBase, filepath.Join(files.HyphaeDir(), "blank.myco"))
	hyphae.Insert(empty)

	oldBase := hyphae.ByName("old").(*hyphae.EmptyHypha)
	old := hyphae.ExtendEmptyToTextual(oldBase, filepath.Join(files.HyphaeDir(), "old.myco"))
	hyphae.Insert(old)
	t.Cleanup(func() {
		hyphae.DeleteHypha(old)
		hyphae.DeleteHypha(empty)
	})

	replace := func(s string) string { return strings.ReplaceAll(s, "old", "new") }
	renameMap, err := renamingPairs([]hyphae.ExistingHypha{old}, replace)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	oldPath := filepath.Join(files.HyphaeDir(), "old.myco")
	newPath := filepath.Join(files.HyphaeDir(), "new.myco")
	if got, ok := renameMap[oldPath]; !ok || got != newPath {
		t.Fatalf("expected mapping %s -> %s, got %v", oldPath, newPath, renameMap)
	}
}
