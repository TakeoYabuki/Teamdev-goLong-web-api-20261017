package service

import (
	"sync"
	"time"

	"github.com/TakeoYabuki/Teamdev-goLong-web-api-20261017/internal/model"
)

const snapshotDuration = 24 * time.Hour

type passSnapshot struct {
	pass    model.Pass
	savedAt time.Time
}

type snapshotStore struct {
	mu        sync.Mutex
	snapshots map[string]passSnapshot
}

func newSnapshotStore() *snapshotStore {
	return &snapshotStore{snapshots: map[string]passSnapshot{}}
}

func snapshotKey(cityID, passID string) string {
	return cityID + "/" + passID
}

func (s *snapshotStore) save(cityID string, passes []model.Pass, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, snapshot := range s.snapshots {
		if now.Sub(snapshot.savedAt) >= snapshotDuration {
			delete(s.snapshots, key)
		}
	}

	for _, pass := range passes {
		s.snapshots[snapshotKey(cityID, pass.ID)] = passSnapshot{pass: pass, savedAt: now}
	}

}

func (s *snapshotStore) get(cityID, passID string, now time.Time) (model.Pass, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	snapshot, found := s.snapshots[snapshotKey(cityID, passID)]
	if !found || now.Sub(snapshot.savedAt) >= snapshotDuration {
		return model.Pass{}, false
	}
	return snapshot.pass, true
}
