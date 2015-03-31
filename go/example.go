// example.go - a test program for the mt19937 package
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

// +build ignore

package main

import (
	"fmt"
	"github.com/funny/mt19937/go"
	"math/rand"
	"time"
)

func main() {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	fmt.Println(rng.NormFloat64())
}
