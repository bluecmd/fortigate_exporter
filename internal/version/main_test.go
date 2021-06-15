// Tests of version parsing
//
// Copyright (C) 2020  Christian Svensson
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package version

import (
	"testing"
)

func TestVersionParseOK(t *testing.T) {
	for _, tv := range []struct {
		v   string
		maj int
		min int
		ok  bool
	}{
		{v: "v6.4.4", maj: 6, min: 4, ok: true},
		{v: "1.0.0", ok: false},
	} {
		t.Run(tv.v, func(t *testing.T) {
			maj, min, ok := ParseVersion(tv.v)
			if !tv.ok {
				if ok {
					t.Errorf("Expected %q to fail to parse, succeeded", tv.v)
				}
				return
			}
			if maj != tv.maj || min != tv.min {
				t.Errorf("Expected %q to be (%d, %d), was (%d, %d)", tv.v, tv.maj, tv.min, maj, min)
			}
		})
	}
}
