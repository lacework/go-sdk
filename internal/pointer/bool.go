package pointer

func BoolPtr(b bool) *bool {
	boolVar := b
	return &boolVar
}

func CompareBoolPtr(ptr *bool, b bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr == b
}
