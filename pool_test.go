// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Pool is no-op under race detector, so all these tests do not work.
//
//go:build !race

package gsync_test

import (
	"testing"

	"github.com/xhofe/gsync"
)

type A struct {
	Name string
}

func TestPool(t *testing.T) {
	// disable GC so we can control when it happens.
	p := gsync.NewPool(func() A {
		return A{}
	})
	a := p.Get()
	a.Name = "a"
	p.Put(a)
	a = p.Get()
	if a.Name != "a" {
		t.Fatalf("got %#v; want a", a)
	}
}
