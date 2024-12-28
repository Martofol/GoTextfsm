package gotextfsm

// Function to check if a slice contains a certain value
func Contains(slice *[]string, value string) bool {
	for _, item := range *slice {
		if item == value {
			return true
		}
	}
	return false
}
