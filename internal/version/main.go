// Parse Fortigate version numbers
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
	"fmt"
)

func ParseVersion(ver string) (int, int, bool) {
	var minor int
	var major int
	n, err := fmt.Sscanf(ver, "v%d.%d.", &major, &minor)
	if err != nil {
		return 0, 0, false
	}
	if n != 2 {
		return 0, 0, false
	}
	return major, minor, true
}
