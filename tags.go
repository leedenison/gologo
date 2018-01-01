package gologo

type ObjectSet map[*Object]bool

func Tag(object *Object, tag string) {
    set, exists := tags[tag]
    if !exists {
        set = make(ObjectSet)
        tags[tag] = set
    }

    set[object] = true
}

func Untag(object *Object, tag string) {
    set, exists := tags[tag]
    if !exists {
        return
    }

    delete(set, object)
}

func UntagAll(object *Object) {
    for _, set := range tags {
        delete(set, object)
    }
}

func HasTag(object *Object, tag string) bool {
    set, exists := tags[tag]
    if !exists {
        return false
    }

    _, exists = set[object]
    return exists
}
