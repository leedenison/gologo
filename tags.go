package gologo

func Tag(object *Object, tag string) {
    set, exists := Tags[tag]
    if !exists {
        set = make(map[*Object]bool)
        Tags[tag] = set
    }

    set[object] = true
}

func Untag(object *Object, tag string) {
    set, exists := Tags[tag]
    if !exists {
        return
    }

    delete(set, object)
}

func UntagAll(object *Object) {
    for _, set := range Tags {
        delete(set, object)
    }
}

func HasTag(object *Object, tag string) bool {
    set, exists := Tags[tag]
    if !exists {
        return false
    }

    _, exists = set[object]
    return exists
}
