package awsutil

// Converts a string -> *string
func StringPointer(s string) *string {
	return &s
}
