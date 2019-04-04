package tokencreatorvalidator

type Validator interface {
	Validate(string) (bool, error)
}

type Creator interface {
	Create() (string, error)
}

type CreatorValidator interface {
	Creator
	Validator
}
