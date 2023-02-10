package slice

func Difference[T comparable](a, b []T) []T {
    mb := make(map[T]struct{}, len(b))
    for _, x := range b {
        mb[x] = struct{}{}
    }
    var diff []T
    for _, x := range a {
        if _, found := mb[x]; !found {
            diff = append(diff, x)
        }
    }
    return diff
}

func ExistsInSlice[T comparable](val T, values []T) bool {
    for _, v := range values {
        if val == v {
            return true
        }
    }
    return false
}