package engine

import (
	"image"
	"io/fs"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
)

// Reflection types used for queries... Is there a better way?
var (
	// TypeOf(pointer to interface).Elem() is "idiomatic" -
	// see https://pkg.go.dev/reflect#example-TypeOf
	AnimerType         = reflect.TypeOf((*Animer)(nil)).Elem()
	BounderType        = reflect.TypeOf((*Bounder)(nil)).Elem()
	ColliderType       = reflect.TypeOf((*Collider)(nil)).Elem()
	DisablerType       = reflect.TypeOf((*Disabler)(nil)).Elem()
	DrawerType         = reflect.TypeOf((*Drawer)(nil)).Elem()
	DrawUpdaterType    = reflect.TypeOf((*DrawUpdater)(nil)).Elem()
	HiderType          = reflect.TypeOf((*Hider)(nil)).Elem()
	IdentifierType     = reflect.TypeOf((*Identifier)(nil)).Elem()
	LoaderType         = reflect.TypeOf((*Loader)(nil)).Elem()
	ParallaxScalerType = reflect.TypeOf((*ParallaxScaler)(nil)).Elem()
	PrepperType        = reflect.TypeOf((*Prepper)(nil)).Elem()
	ScannerType        = reflect.TypeOf((*Scanner)(nil)).Elem()
	ScenerType         = reflect.TypeOf((*Scener)(nil)).Elem()
	SaverType          = reflect.TypeOf((*Saver)(nil)).Elem()
	TransformerType    = reflect.TypeOf((*Transformer)(nil)).Elem()
	UpdaterType        = reflect.TypeOf((*Updater)(nil)).Elem()

	// Behaviours lists all the behaviours that can be queried with Game.Query.
	Behaviours = []reflect.Type{
		AnimerType,
		BounderType,
		ColliderType,
		DisablerType,
		DrawerType,
		DrawUpdaterType,
		HiderType,
		IdentifierType,
		LoaderType,
		ParallaxScalerType,
		PrepperType,
		ScannerType,
		ScenerType,
		SaverType,
		TransformerType,
		UpdaterType,
	}
)

// Animer components have a current frame index.
type Animer interface {
	Updater

	CurrentFrame() int
	Reset()
}

// Bounder components have a bounding rectangle.
type Bounder interface {
	BoundingRect() image.Rectangle
}

// Collider components have tangible form.
type Collider interface {
	CollidesWith(image.Rectangle) bool
}

// Disabler components can be disabled.
type Disabler interface {
	IsDisabled() bool
	Disable()
	Enable()
}

// Drawer components can draw themselves. Draw is called often. Each component
// must call Draw on any internal components not known to the engine (i.e. not
// passed to Game.Register or returned from Scan).
type Drawer interface {
	Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions)
	DrawOrder() float64
}

// DrawUpdater components can be both drawn and updated.
// Same comments as for Drawer and Updater.
type DrawUpdater interface {
	Drawer
	Updater
}

// Hider components can be hidden.
type Hider interface {
	IsHidden() bool
	Hide()
	Show()
}

// Identifier components have a sense of self. This makes it easier for
// components to find and interact with one another.
type Identifier interface {
	Ident() string
}

// Loader components get the chance to load themselves. This happens
// before preparation.
type Loader interface {
	Load(fs.FS) error
}

// ParallaxScaler components have a scaling factor. This is used for
// parallax layers in a scene, and can be thought of as 1/distance.
type ParallaxScaler interface {
	ParallaxFactor() float64
}

// Prepper components can be prepared. It is called after the component
// database has been populated but before the game is run. The component can
// store the reference to game, if needed, and also query the component database.
type Prepper interface {
	Prepare(game *Game) error
}

// Scanner components can be scanned. It is called when the game tree is walked
// (such as when the game component database is constructed).
// Scan should return a slice containing all immediate subcomponents.
type Scanner interface {
	Scan() []interface{}
}

// Scener components are a scene (Scene or SceneRef).
type Scener interface {
	// Q: Why not make Scene able to load itself?
	// A: Having separate types makes it easier to reason about what is loading
	//    what. There is less ambiguity about what "save" means (the contents of
	//    the scene, or the path to the file?) Additionally, the gob decoder
	//    decodes over existing fields, which could lead to some fun bugs.
	//
	// Q: Why not make Scener a small interface, e.g. with just Scene() ?
	// A: Everything in the engine would then need to type switch for Scener or
	//    SceneRef, i.e.
	//        switch x := i.(type) {
	//        case Drawer:
	//            i.Draw(screen, opts)
	//        case Scener:
	//            i.Scene().Draw(screen, opts)
	//        }
	//    It seems cleaner to let the engine assert only for the interface it
	//    needs at that moment.

	Bounder
	Disabler
	Hider
	Identifier
	Scanner
	Transformer
}

// Saver components can be saved to disk.
type Saver interface {
	Save() error
}

// Transformer components can transform their child components.
type Transformer interface {
	Transform() ebiten.DrawImageOptions
}

// Updater components can update themselves. Update is called repeatedly. Each
// component must call Update on any internal components not known to the engine
//  (i.e. not passed to Game.Register or returned from Scan).
type Updater interface {
	Update() error
}
