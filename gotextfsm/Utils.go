package gotextfsm

import "log"

// Function to check if a slice contains a certain value
func Contains(slice []string, value string) bool {
	log.Println("I m inside Contains function Sended Value :", value)
	for _, item := range slice {
		log.Println("Inside For Looop The Item is", item)
		if item == value {
			return true
		}
	}
	return false
}
