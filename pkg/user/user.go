package user

// Status of user
type Status struct {
	LoggedIn bool
	LastTime int
	LastIP   int
}

// GetStatus - returns the login status of the user
func GetStatus() Status {
	return Status{}
}

// LogOn to the system
func LogOn(email string, password string) {

}

// LogOff from the system
func LogOff() {

}

// Register with the system
func Register(name string, email string, password string) {

}
