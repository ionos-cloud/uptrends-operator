package finalizers

// Finalizer is an interface that can be used to add, remove and check for
type Finalizer interface {
	GetFinalizers() []string
}

// AddFinalizer ...
func AddFinalizer(obj Finalizer, finalizer string) []string {
	finalizers := obj.GetFinalizers()
	for _, f := range finalizers {
		if finalizer == f {
			return finalizers
		}
	}

	return append([]string{finalizer}, finalizers...)
}

// HasFinalizer ...
func HasFinalizer(obj Finalizer, finalizer string) bool {
	finalizers := obj.GetFinalizers()
	for _, f := range finalizers {
		if finalizer == f {
			return true
		}
	}
	return false
}

// RemoveFinalizer ...
func RemoveFinalizer(obj Finalizer, finalizer string) []string {
	finalizers := []string{}
	for _, f := range obj.GetFinalizers() {
		if f == finalizer {
			continue
		}
		finalizers = append(finalizers, f)
	}
	return finalizers
}
