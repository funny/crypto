// doc.go - package documentation for github.com/seehuhn/mt19937
// Copyright (C) 2013  Jochen Voss <voss@seehuhn.de>
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

// Package mt19937 is a pure-go implementation of the 64bit Mersenne
// Twister pseudo random number generator (PRNG).  The Mersenne
// Twister, developed by Takuji Nishimura and Makoto Matsumoto, is,
// for example, commonly used in Monte Carlo simulations.  The
// implementation in the mt19937 package closely follows the reference
// implementation from
//
//     http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/emt64.html
//
// and, for identical seeds, gives identical output to the reference
// implementation.
//
// The PRNG from the mt19937 package implements the rand.Source
// interface from the math/rand package.  Typically the PRNG is
// wrapped in a rand.Rand object as in the following example:
//
//     rng := rand.New(mt19937.New())
//     rng.Seed(time.Now().UnixNano())
package mt19937
