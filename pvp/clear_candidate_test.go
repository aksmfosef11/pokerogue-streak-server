// SPDX-License-Identifier: AGPL-3.0-or-later
package pvp

import (
	"reflect"
	"testing"

	"github.com/pagefaultgames/rogueserver/defs"
)

func TestRecoverClearCandidate(t *testing.T) {
	const version = "1.12.2.0"
	candidate := defs.SessionSaveData{Seed: "issued-seed", GameVersion: version, GameMode: 0, WaveIndex: 200, BattleType: 2, Money: 100}
	makeState := func() *State {
		return &State{Runs: map[string]*Run{"run": {ID: "run", Seed: candidate.Seed, Version: version, TrustEpoch: 3, Anchors: []Anchor{{Wave: 200}}}}, TrustEpoch: 3}
	}

	t.Run("records once and never overwrites", func(t *testing.T) {
		st := makeState()
		if err := recoverClearCandidate(st, "run", version, candidate, 123); err != nil {
			t.Fatal(err)
		}
		first := *st.Runs["run"].Candidate
		if err := recoverClearCandidate(st, "run", version, candidate, 124); err != nil {
			t.Fatal(err)
		}
		other := candidate
		other.Money = 999
		if err := recoverClearCandidate(st, "run", version, other, 125); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(*st.Runs["run"].Candidate, first) || st.Runs["run"].ClearedAt != 123 {
			t.Fatal("retry changed the original candidate or clear time")
		}
	})

	for _, tc := range []struct {
		name   string
		mutate func(*State, *defs.SessionSaveData)
	}{
		{"missing run", func(s *State, _ *defs.SessionSaveData) { delete(s.Runs, "run") }},
		{"run ID mismatch", func(s *State, _ *defs.SessionSaveData) { s.Runs["run"].ID = "different" }},
		{"seed mismatch", func(_ *State, c *defs.SessionSaveData) { c.Seed = "other" }},
		{"version mismatch", func(_ *State, c *defs.SessionSaveData) { c.GameVersion = "1.12.1.3" }},
		{"run version mismatch", func(s *State, _ *defs.SessionSaveData) { s.Runs["run"].Version = "1.12.1.3" }},
		{"not classic", func(_ *State, c *defs.SessionSaveData) { c.GameMode = 1 }},
		{"wrong wave", func(_ *State, c *defs.SessionSaveData) { c.WaveIndex = 199 }},
		{"not clear", func(_ *State, c *defs.SessionSaveData) { c.BattleType = 0 }},
		{"missing 200 anchor", func(s *State, _ *defs.SessionSaveData) { s.Runs["run"].Anchors = []Anchor{{Wave: 199}} }},
		{"reset epoch", func(s *State, _ *defs.SessionSaveData) { s.TrustEpoch++ }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			st, c := makeState(), candidate
			tc.mutate(st, &c)
			before := st.Runs["run"]
			if err := recoverClearCandidate(st, "run", version, c, 123); err == nil {
				t.Fatal("invalid candidate accepted")
			}
			if before != nil && (before.Candidate != nil || before.ClearedAt != 0) {
				t.Fatal("invalid candidate mutated run")
			}
		})
	}
}
