package engine

import (
	"testing"

	"drjosh.dev/gurgle/geom"
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
