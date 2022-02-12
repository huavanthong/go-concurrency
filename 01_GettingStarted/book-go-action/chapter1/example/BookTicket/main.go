/***
	https://www.meisternote.com/app/note/sZFeNENApMk2/book-ticket-application
***/
package main

// Package sync provides basic synchroization primitives such as mutual exclusion locks.
// Other than the Once and WaitGroup types, most are intended for use by low-level library routines.
// Higher-level synchronization is better done via channels and communication.
import (
	"fmt"
	"sync"
	"time"
)

// Giống như C, đây là cách để khai báo một constant value, hoặc các variable nằm ở global variable area
const conferenceTickets int = 50

var conferenceName = "Go Conference"
var remainingTickets uint = 50

// The slice bookings refer to a new array of 100 ints
var bookings = make([]UserData, 0)

// Create a structure for UserData
type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
}

// Create a Wait group to wait for a collection of goroutines to finish
var wg = sync.WaitGroup{}

func main() {

	greetUsers()

	for {
		firstName, lastName, email, userTickets := getUserInput()
		// Ngoài ra, ta có thể sử dụng các file bên ngoài như: helper.go để xây dựng các API cho source này.
		isValidName, isValidEmail, isValidTicketNumber := validateUserInput(firstName, lastName, email, userTickets)

		if isValidName && isValidEmail && isValidTicketNumber {
			// Nếu các các info của user được validate, thì ta sẽ tiến hành book Ticket.
			bookTicket(userTickets, firstName, lastName, email)
			// Increment the WaitGroup counter.
			wg.Add(1)
			// Create a Goroutine để bắt đầu cho việc send result booking
			// Problem:
			// Têm một problem nữa nằm ở đâu, tại sao ở đây ta không thấy việc tạo Channel mà lại có thể sử
			// dụng được go().
			// Solution:
			// https://www.meisternote.com/app/note/sZFeNENApMk2/book-ticket-application
			go sendTicket(userTickets, firstName, lastName, email)

			firstNames := getFirstNames()
			fmt.Printf("The first names of bookings are: %v\n", firstNames)

			if remainingTickets == 0 {
				// end program
				fmt.Println("Our conference is booked out. Come back next year.")
				// break
			}
		} else {
			if !isValidName {
				fmt.Println("first name or last name you entered is too short")
			}
			if !isValidEmail {
				fmt.Println("email address you entered doesn't contain @ sign")
			}
			if !isValidTicketNumber {
				fmt.Println("number of tickets you entered is invalid")
			}
		}
	}
	wg.Wait()
}

func greetUsers() {
	fmt.Printf("Welcome to %v booking application\n", conferenceName)
	fmt.Printf("We have total of %v tickets and %v are still available.\n", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend")
}

func getFirstNames() []string {
	firstNames := []string{}
	for _, booking := range bookings {
		firstNames = append(firstNames, booking.firstName)
	}
	return firstNames
}

// Để xây dựng được một function có thể return nhiều result parameter ta sử dụng feature #NameResultParameters
// https://www.meisternote.com/app/note/VWIHKl5VBJK8/named-result-parameters
// Với syntax cho cách 1 như sau:
// func function_name(Parameter-list)(result_parameter1 data-_type, result_parameter2 data_type, ….){ ... }
// Với syntax cho cách 2 như sau:
// func function_name(Parameter-list)(data-_type, data_type, ….){ ... }
func getUserInput() (string, string, string, uint) {

	// Nhưng ta thấy rằng, (string, string, string, uint) nó không có define result_name, nhưng sẽ được return ở cuối function.
	// Vậy đây là cách thứ 2 khi sử dụng NameResultParameters
	var firstName string
	var lastName string
	var email string
	var userTickets uint

	fmt.Println("Enter your first name: ")
	fmt.Scan(&firstName)

	fmt.Println("Enter your last name: ")
	fmt.Scan(&lastName)

	fmt.Println("Enter your email address: ")
	fmt.Scan(&email)

	fmt.Println("Enter number of tickets: ")
	fmt.Scan(&userTickets)

	return firstName, lastName, email, userTickets
}

func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	// Việc các client gửi request để book ticket, và việc check số vé còn lại phải được làm ở đây. Ngay khi request bắt đầu.
	remainingTickets = remainingTickets - userTickets

	var userData = UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}

	// Bởi vì ta đã tạo Slice từ đầu program, nên cứ mỗi booking ok thì ta cứ Append() vào Database của ta.
	bookings = append(bookings, userData)
	// Trong Golang, ta có quyển sử dụng %v format để hiển thị cả structure của info.
	fmt.Printf("List of bookings is %v\n", bookings)

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", remainingTickets, conferenceName)
}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {

	// Giả sử việc gửi các Ticket cần set up từa lưa thứ, và nó tốn khoảng 50 giây.
	// Và problem là ở đây:
	// Nếu việc gửi result tốn rất nhiều thời gian, nhưng ta có 2 kết quả có thể xảy ra tại đây:
	// 1. Là nếu ticket được book OK, thì số vé còn lại phải bị giảm.
	// 2. Giả sử như ticket không book được, thì số vé còn lại phải được recover.
	// Làm sao giải quyết được vấn đề ở đây?
	time.Sleep(50 * time.Second)
	var ticket = fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
	fmt.Println("#################")
	fmt.Printf("Sending ticket:\n %v \nto email address %v\n", ticket, email)
	fmt.Println("#################")
	wg.Done()
}
