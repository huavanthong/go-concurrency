package main

import "strings"

// Bởi vì chữ cái đầu tiên không được viết hoa, nên đây không phải là external API cho các file cùng package được sử dụng.
// Vậy để sử dụng được function này, thì ta phải build file này trong Context của project.
func validateUserInput(firstName string, lastName string, email string, userTickets uint) (bool, bool, bool) {
	isValidName := len(firstName) >= 2 && len(lastName) >= 2
	isValidEmail := strings.Contains(email, "@")
	isValidTicketNumber := userTickets > 0 && userTickets <= remainingTickets
	return isValidName, isValidEmail, isValidTicketNumber
}
