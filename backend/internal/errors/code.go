package errors

const (
	// Generic
	BadRequest = "bad_request" // Request body is malformed
	Internal   = "internal"    // Internal error of some sort
	Invalid    = "invalid"     // Request is invalid
	NotFound   = "not_found"   // Resource not found
	Unexpected = "unexpected"  // Unexpected error (when the error has no code)

	// Auth
	ChallengeAlreadyUsed = "challenge_already_used" // The challenge has already been verified
	ChallengeExpired     = "challenge_expired"      // The challenge has expired
	ChallengeFailed      = "challenge_failed"       // The secret provided was incorrect
	ChallengeInvalidated = "challenge_invalidated"  // The challenge has been flagged as suspicious
	ChallengeMaxAttempts = "challenge_max_attempts" // The maximum number of attempts has been reached
	DoNotContact         = "do_not_contact"         // Email has been added to the do-not-contact list
	EmailFlagged         = "email_flagged"          // Email temporarily flagged
	EmailInvalid         = "email_invalid"          // Email is unrecognized
	EmailUnavailable     = "email_unavailable"      // Email is already in use
	IPFlagged            = "ip_flagged"             // IP address temporarily flagged
	SessionExpired       = "session_expired"        // Session has expired
	SessionInvalid       = "session_invalid"        // Session token is unrecognized
	Unauthenticated      = "unauthenticated"        // No session cookie provided

	// Access
	Forbidden = "forbidden" // Access to this API is denied based on our access policies
)

const (
	internalMessage  = "an internal error occurred"
	forbiddenMessage = "forbidden"
	notFoundMessage  = "not found"
)

var codeToStatus = map[string]int{
	BadRequest:           400,
	ChallengeAlreadyUsed: 422,
	ChallengeExpired:     422,
	ChallengeFailed:      422,
	ChallengeInvalidated: 422,
	ChallengeMaxAttempts: 422,
	DoNotContact:         422,
	EmailFlagged:         422,
	EmailInvalid:         422,
	EmailUnavailable:     422,
	Forbidden:            403,
	Internal:             500,
	Invalid:              422,
	IPFlagged:            422,
	NotFound:             404,
	SessionExpired:       401,
	SessionInvalid:       401,
	Unauthenticated:      401,
	Unexpected:           500,
}
