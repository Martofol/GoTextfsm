package main

import "fmt"

func removeDuplicates(strings []string) []string {
	// Create a map to store unique strings
	uniqueStrings := make(map[string]bool)
	var result []string

	// Loop through the slice and add unique strings to the result
	for _, str := range strings {
		if _, found := uniqueStrings[str]; !found {
			uniqueStrings[str] = true
			result = append(result, str)
		}
	}

	return result
}

func main() {
	// Example string slice with duplicates
	strings := []string{"apple", "banana", "apple", "orange", "banana", "pear"}

	// Remove duplicates
	uniqueStrings := removeDuplicates(strings)

	// Print the result
	fmt.Println(uniqueStrings)
}
