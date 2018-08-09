package gologo

// ObjectSet : struct to hold the objects for each tag
type ObjectSet map[*Object]bool

// Tag : Tag the supplied object in the tags object list
func Tag(object *Object, tag string) {
	set, exists := tags[tag]
	if !exists {
		set = make(ObjectSet)
		tags[tag] = set
	}

	set[object] = true
}

// Untag : Untag the supplied object in the tags object list
func Untag(object *Object, tag string) {
	set, exists := tags[tag]
	if !exists {
		return
	}

	delete(set, object)
}

// UntagAll : Remove object from all tags object lists
func UntagAll(object *Object) {
	for _, set := range tags {
		delete(set, object)
	}
}

// HasTag : Returns true if the specified tag's object list contains object
func HasTag(object *Object, tag string) bool {
	set, exists := tags[tag]
	if !exists {
		return false
	}

	_, exists = set[object]
	return exists
}
