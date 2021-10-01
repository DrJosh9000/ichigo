/*
Copyright 2021 Josh Deprez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package engine

import (
	"testing"

	"github.com/DrJosh9000/ichigo/geom"
	"github.com/google/go-cmp/cmp"
	"github.com/hajimehoshi/ebiten/v2"
)

type fakeDrawBoxer string

func (fakeDrawBoxer) Draw(*ebiten.Image, *ebiten.DrawImageOptions) {}
func (fakeDrawBoxer) BoundingBox() geom.Box {
	return geom.Box{}
}

func TestDAGAddRemoveEdge(t *testing.T) {
	d := make(dag)
	var u, v fakeDrawBoxer = "u", "v"
	d.addEdge(u, v)
	want := dag{
		u: edges{
			in:  drawerSet{},
			out: drawerSet{v: {}},
		},
		v: edges{
			in:  drawerSet{u: {}},
			out: drawerSet{},
		},
	}
	if diff := cmp.Diff(d, want, cmp.AllowUnexported(edges{})); diff != "" {
		t.Errorf("diff after adding edge (u->v):\n%s", diff)
	}
	d.removeEdge(u, v)
	want = dag{
		u: edges{
			in:  drawerSet{},
			out: drawerSet{},
		},
		v: edges{
			in:  drawerSet{},
			out: drawerSet{},
		},
	}
	if diff := cmp.Diff(d, want, cmp.AllowUnexported(edges{})); diff != "" {
		t.Errorf("diff after removing edge (u->v):\n%s", diff)
	}
}

func TestDAGAddVertex(t *testing.T) {
	d := make(dag)
	u := fakeDrawBoxer("u")
	d.addVertex(u)
	want := dag{
		u: edges{
			in:  drawerSet{},
			out: drawerSet{},
		},
	}
	if diff := cmp.Diff(d, want, cmp.AllowUnexported(edges{})); diff != "" {
		t.Errorf("diff after adding vertex u:\n%s", diff)
	}
	v := fakeDrawBoxer("v")
	d.addVertex(v)
	want = dag{
		u: edges{
			in:  drawerSet{},
			out: drawerSet{},
		},
		v: edges{
			in:  drawerSet{},
			out: drawerSet{},
		},
	}
	if diff := cmp.Diff(d, want, cmp.AllowUnexported(edges{})); diff != "" {
		t.Errorf("diff after adding vertex v:\n%s", diff)
	}
	d.addVertex(u)
	if diff := cmp.Diff(d, want, cmp.AllowUnexported(edges{})); diff != "" {
		t.Errorf("diff after adding vertex u a second time:\n%s", diff)
	}
}

func TestDAGRemoveVertex(t *testing.T) {
	u := fakeDrawBoxer("u")
	v := fakeDrawBoxer("v")
	w := fakeDrawBoxer("w")
	d := dag{
		u: edges{
			in:  drawerSet{},
			out: drawerSet{v: {}},
		},
		v: edges{
			in:  drawerSet{u: {}},
			out: drawerSet{w: {}},
		},
		w: edges{
			in:  drawerSet{v: {}},
			out: drawerSet{},
		},
	}
	d.removeVertex(u)
	want := dag{
		v: edges{
			in:  drawerSet{},
			out: drawerSet{w: {}},
		},
		w: edges{
			in:  drawerSet{v: {}},
			out: drawerSet{},
		},
	}
	if diff := cmp.Diff(d, want, cmp.AllowUnexported(edges{})); diff != "" {
		t.Errorf("diff after removing vertex u:\n%s", diff)
	}
	d.removeVertex(u)
	if diff := cmp.Diff(d, want, cmp.AllowUnexported(edges{})); diff != "" {
		t.Errorf("diff after removing vertex u a second time:\n%s", diff)
	}
	d.removeVertex(w)
	want = dag{
		v: edges{
			in:  drawerSet{},
			out: drawerSet{},
		},
	}
	if diff := cmp.Diff(d, want, cmp.AllowUnexported(edges{})); diff != "" {
		t.Errorf("diff after removing vertex v:\n%s", diff)
	}
}

func TestTopWalk(t *testing.T) {
	u := fakeDrawBoxer("u")
	v := fakeDrawBoxer("v")
	w := fakeDrawBoxer("w")
	d := dag{
		u: edges{
			in:  drawerSet{},
			out: drawerSet{v: {}},
		},
		v: edges{
			in:  drawerSet{u: {}},
			out: drawerSet{w: {}},
		},
		w: edges{
			in:  drawerSet{v: {}},
			out: drawerSet{},
		},
	}

	var got []Drawer
	d.topWalk(func(x Drawer) {
		got = append(got, x)
	})
	want := []Drawer{u, v, w}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("topWalk visited in wrong order - diff:\n%s", diff)
	}
}

func TestTopWalkToleratesCycle(t *testing.T) {
	u := fakeDrawBoxer("u")
	v := fakeDrawBoxer("v")
	w := fakeDrawBoxer("w")
	d := dag{
		u: edges{
			in:  drawerSet{w: {}},
			out: drawerSet{v: {}},
		},
		v: edges{
			in:  drawerSet{u: {}},
			out: drawerSet{w: {}},
		},
		w: edges{
			in:  drawerSet{v: {}},
			out: drawerSet{u: {}},
		},
	}
	got := make(map[Drawer]int)
	d.topWalk(func(x Drawer) {
		got[x]++
	})
	want := map[Drawer]int{u: 1, v: 1, w: 1}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("topWalk visited vertices wrong number of times - diff:\n%s", diff)
	}
}
