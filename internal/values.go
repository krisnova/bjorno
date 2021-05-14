package internal

// TODO @kris-nova we should make this dynamic
//
// We have an opportunity here to make this a dynamic server
// and allow users to pass in their own Values struct and their
// own refresh methods.

// Values is the droid you are looking for.
//
// If this struct has data, it will interpolated
// at runtime using Go's text/template.
//
// Web development is hard.
type Values struct {
	Title  string
	Author string
	Beeps  string
}

var v *Values

// GetValues is the magic sauce for where we get our
// values from.
//
// Here we always quickly return as this will be
// pulling at runtime. The guarantee we have with
// any of these Values is "best effort".
//
// As we are waiting on data, if it's not there
// we simply don't return it.
//
// However we let it try to "catch up" concurrently
// in the hopes we get whatever it is we are looking for.
//
// Note: This is where the DoS attacks will kill us.
func GetValues() *Values {
	go RefreshValues()
	return v
}

// RefreshValues is a procederual method to "refresh"
// whatever values we have in the cache.
func RefreshValues() error {
	v = &Values{
		Title:  "nivenly.com",
		Author: "kris nóva",
		Beeps:  "boops",
	}
	return nil
}
