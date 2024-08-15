package pointer

// CompareBoolPtr compares a bool pointer to a bool primitive
// Returns false if bool pointer is nil
func CompareBoolPtr(ptr *bool, b bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr == b
}
