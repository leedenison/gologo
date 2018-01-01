package gologo

func containsInt(s []int, v int) bool {
    for _, c := range s {
        if c == v {
            return true
        }
    }
    return false
}

