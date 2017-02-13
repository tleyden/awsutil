package awsutil

// StringPointer converts a string -> *string
func StringPointer(s string) *string {
	return &s
}

// Int64Pointer converts an int64 -> *int64
func Int64Pointer(i int64) *int64 {
	return &i
}