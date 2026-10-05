package runmgr

import (
	"github.com/mohammad-safakhou/diffmind/internal/workspace/store"
	"testing"
)

func TestProjectPackStorageKeysResolveManifestIDs(t *testing.T) {
	m, pid, refs := setup(t)
	key, err := m.store.CreatePack(pid, "example.conventions", []byte(`{"id":"example.conventions","name":"Example"}`))
	if err != nil {
		t.Fatal(err)
	}
	for _, selector := range []string{key, "example.conventions"} {
		if _, err = m.store.UpdateRepo(pid, refs[0].RepoID, func(r *store.Repo) { r.PackIDs = []string{selector} }); err != nil {
			t.Fatal(err)
		}
		cfg, _, err := m.buildConfig(pid, store.RunManifest{Repos: refs})
		if err != nil {
			t.Fatal(err)
		}
		if len(cfg.Repos.ServiceRepos[0].PackIDs) != 1 || cfg.Repos.ServiceRepos[0].PackIDs[0] != "example.conventions" {
			t.Fatalf("selection omitted: %+v", cfg.Repos.ServiceRepos[0])
		}
	}
	if err = m.store.DeletePack(pid, key); err != nil {
		t.Fatal(err)
	}
	if _, _, err = m.buildConfig(pid, store.RunManifest{Repos: refs}); err == nil {
		t.Fatal("missing selected pack silently omitted")
	}
}
