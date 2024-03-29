package tags

import (
	"github.com/leedenison/gologo"
)

// ObjectSet : struct to hold the objects for each tag
type ObjectSet map[*gologo.Object]bool

// TagSet : maps string tags to sets of objects
type TagSet map[string]ObjectSet

// Tag : Tag the supplied object in the tags object list
func (t TagSet) Tag(object *gologo.Object, tag string) {
	set, exists := t[tag]
	if !exists {
		set = make(ObjectSet)
		t[tag] = set
	}

	set[object] = true
}

// Untag : Untag the supplied object in the tags object list
func (t TagSet) Untag(object *gologo.Object, tag string) {
	set, exists := t[tag]
	if !exists {
		return
	}

	delete(set, object)
}

// UntagAll : Remove object from all tags object lists
func (t TagSet) UntagAll(object *gologo.Object) {
	for _, set := range t {
		delete(set, object)
	}
}

// HasTag : Returns true if the specified tag's object list contains object
func (t TagSet) HasTag(object *gologo.Object, tag string) bool {
	set, exists := t[tag]
	if !exists {
		return false
	}

	_, exists = set[object]
	return exists
}

// GetAll : Returns a slice of all the objects tagged with tag
func (t TagSet) GetAll(tag string) []*gologo.Object {
	set, exists := t[tag]
	if !exists {
		return []*gologo.Object{}
	}

	keys := make([]*gologo.Object, len(set))

	i := 0
	for k := range set {
		keys[i] = k
		i++
	}

	return keys
}
