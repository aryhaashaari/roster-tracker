package rule

import "regexp"

var (
	alphabet                        = regexp.MustCompile(`^[A-Za-z]+$`)
	numeric                         = regexp.MustCompile(`^[0-9]+$`)
	alphabetNumeric                 = regexp.MustCompile(`^[^,'"!@#$%&*()?\n\\/{}]*$`)
	alphabetNumericSpecialCharacter = regexp.MustCompile(`^[\w\s\(\)\-\+\,\.\!\?\/\\]+$`)
	validUUID                       = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	alphabetContainsUnderscore      = regexp.MustCompile(`^[a-zA-Z]+(_[a-zA-Z]+)?$`)
	alphaNumericUnderscore          = regexp.MustCompile(`^[\w]+$`)
	notAllowXXS                     = regexp.MustCompile(`^[^<>&]+$`)
)
