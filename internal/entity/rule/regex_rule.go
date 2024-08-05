package rule

import validation "github.com/go-ozzo/ozzo-validation/v4"

var (
	AlphabetNumericSpecialCharacterRule = validation.Match(alphabetNumericSpecialCharacter).Error(`must be among or combination these characters (a-z, A-Z, 0-9, space, enter, tab, comma(,), dot(.), slash and back slash(\/), question mark(?), parentheses('()'),exclamation mark(!), underscore(_), plus and minus(-+))`)
	ValidUUIDRule                       = validation.Match(validUUID).Error("must be valid format UUID")
	AlphabetNumericRule                 = validation.Match(alphabetNumeric).Error("invalid format")
	AlphabetContainsUnderscoreRule      = validation.Match(alphabetContainsUnderscore).Error("invalid format")
	AlphaNumericUnderscoreRule          = validation.Match(alphaNumericUnderscore).Error(`must be among or combination these characters (a-z, A-Z, 0-9, and underscore(_))`)
	AlphabetRule                        = validation.Match(alphabet).Error(`must be among or combination these characters (a-z, A-Z)`)
	NumericRule                         = validation.Match(numeric).Error("must be numeric")
	NowAllowXXSRule                     = validation.Match(notAllowXXS).Error("invalid input format")
)
