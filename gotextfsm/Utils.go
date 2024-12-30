package gotextfsm

import (
	"log"
	"regexp"
)

// Function to check if a slice contains a certain value
func Contains(slice *[]string, value string) bool {
	for _, item := range *slice {
		if item == value {
			return true
		}
	}
	return false
}

func GetNamedMatches(r *regexp.Regexp, s string) map[string]string {
	match := r.FindStringSubmatch(s)
	if match == nil {
		return nil
	}
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	log.Println("regrex = ", r)
	log.Println("Match = ", match)
	log.Println("subMatch = ", subMatchMap)
	return subMatchMap
}
