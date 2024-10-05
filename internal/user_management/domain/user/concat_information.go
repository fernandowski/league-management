package domain

type ContactInformation struct {
	firstName string
	lastName  string
}

func NewContactInformation(firstName string, lastName string) *ContactInformation {
	return &ContactInformation{firstName: firstName, lastName: lastName}
}
