package apperror

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code    string
	Message string
}

// Error implements the error interface for AppError
func (e AppError) Error() string {
	return e.Message
}

// Code extracts the error Code from an error, returning the AppError Code if it's an AppError, otherwise returns INTERNAL_ERROR
func Code(err error) string {
	if isAppError(err) {
		return err.(AppError).Code
	}
	return ErrInternal.Code
}
func NewError(originalErr error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Message extracts the error Message from an error, returning the AppError Message if it's an AppError, otherwise returns a generic internal error Message
func Message(err error) string {
	if isAppError(err) {
		return err.(AppError).Message
	}
	return ErrInternal.Message
}

// isAppError checks if an error is an AppError (safe to expose to frontend)
func isAppError(err error) bool {
	var appError AppError
	ok := errors.As(err, &appError)
	return ok
}

// isErrorType checks if err matches any of the provided target errors
func isErrorType(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// StatusFromError maps custom errors to HTTP status codes
func StatusFromError(err error) int {
	switch {
	// 400 Bad Request
	case isErrorType(err, ErrBadRequest, ErrInvalidID, ErrInvalidMembershipData, ErrInvalidOTP, ErrOTPExpired):
		return http.StatusBadRequest
	// 401 Unauthorized
	case isErrorType(err, ErrInvalidCredentials, ErrInvalidToken, ErrInvalidClaims, ErrInvalidIssuer, ErrInvalidAudience, ErrTokenInvalidated):
		return http.StatusUnauthorized
	// 403 Forbidden
	case isErrorType(err, ErrForbidden, ErrUserInactive, ErrEmailNotVerified):
		return http.StatusForbidden
	// 404 Not Found
	case isErrorType(err, ErrUserNotFound, ErrMembershipNotFound, ErrInvitationNotFound, ErrClassroomNotFound, ErrJoinRequestNotFound, ErrPostNotFound):
		return http.StatusNotFound
	// 409 Conflict
	case isErrorType(err, ErrUsernameExists, ErrEmailExists, ErrAlreadyMember, ErrEmailAlreadyVerified, ErrLoginMethodMismatch, ErrInvitationAlreadyExists, ErrAlreadyInClassroom, ErrInvitationExpired, ErrInvitationAlreadyProcessed, ErrStudentCodeAlreadyUsed, ErrJoinRequestAlreadyProcessed):
		return http.StatusConflict
	// 500 Internal Server Error
	case isErrorType(err, ErrInternal, ErrNoFieldsToUpdate, ErrMembershipCreateFailed, ErrMembershipDeleteFailed):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

var (
	// Auth-related
	ErrInvalidCredentials   = AppError{Code: "INVALID_CREDENTIALS", Message: "Invalid username or password"}
	ErrInvalidToken         = AppError{Code: "INVALID_TOKEN", Message: "Invalid or expired token"}
	ErrInvalidClaims        = AppError{Code: "INVALID_CLAIMS", Message: "Invalid token claims"}
	ErrInvalidIssuer        = AppError{Code: "INVALID_ISSUER", Message: "Invalid token issuer"}
	ErrInvalidAudience      = AppError{Code: "INVALID_AUDIENCE", Message: "Invalid token audience"}
	ErrTokenInvalidated     = AppError{Code: "TOKEN_INVALIDATED", Message: "Token has been invalidated"}
	ErrForbidden            = AppError{Code: "FORBIDDEN", Message: "You do not have permission to perform this action"}
	ErrBadRequest           = AppError{Code: "BAD_REQUEST", Message: "Bad request"}
	ErrEmailNotVerified     = AppError{Code: "EMAIL_NOT_VERIFIED", Message: "Email has not been verified"}
	ErrEmailAlreadyVerified = AppError{Code: "EMAIL_ALREADY_VERIFIED", Message: "Email has already been verified"}
	ErrInvalidOTP           = AppError{Code: "INVALID_OTP", Message: "Invalid verification code"}
	ErrOTPExpired           = AppError{Code: "OTP_EXPIRED", Message: "Verification code has expired"}
	ErrLoginMethodMismatch  = AppError{Code: "LOGIN_METHOD_MISMATCH", Message: "This email is registered with a different login method. Please use the original method."}
	ErrEmailNotRegistered   = AppError{Code: "EMAIL_NOT_REGISTERED", Message: "Email chưa được đăng ký"}

	// Generic
	ErrInternal          = AppError{Code: "INTERNAL_ERROR", Message: "Internal server error"}
	ErrNoFieldsToUpdate  = AppError{Code: "NO_FIELDS_TO_UPDATE", Message: "No fields provided to update"}
	ErrInvalidID         = AppError{Code: "INVALID_ID", Message: "Invalid ID format"}
	ErrPaginationInvalid = AppError{Code: "PAGINATION_INVALID", Message: "Page number or page size is invalid. Page size must be smaller than 500."}
	ErrNotFound          = AppError{Code: "NOT_FOUND", Message: "Resource not found"}

	// User-related
	ErrUserNotFound   = AppError{Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrUsernameExists = AppError{Code: "USERNAME_EXISTS", Message: "Username already exists"}
	ErrEmailExists    = AppError{Code: "EMAIL_EXISTS", Message: "Email already exists"}
	ErrUserInactive   = AppError{Code: "USER_INACTIVE", Message: "User account is inactive"}

	// Group-related
	ErrGroupNotFound        = AppError{Code: "GROUP_NOT_FOUND", Message: "Group not found"}
	ErrProjectGroupNotFound = AppError{Code: "PROJECT_GROUP_NOT_FOUND", Message: "Project group not found"}
	ErrProjectNotFound      = AppError{Code: "PROJECT_NOT_FOUND", Message: "Project not found"}
	ErrInvalidMemberNumber  = AppError{Code: "INVALID_MEMBER_NUMBER", Message: "Invalid member number"}

	// Project-related
	ErrReportAlreadyExists      = AppError{Code: "REPORT_ALREADY_EXISTS", Message: "Report already exists for this period"}
	ErrReportPeriodNotFound     = AppError{Code: "REPORT_PERIOD_NOT_FOUND", Message: "Report period not found"}
	ErrProjectRoundNotFound     = AppError{Code: "PROJECT_ROUND_NOT_FOUND", Message: "Project round not found"}
	ErrProjectGroupLimitReached = AppError{Code: "PROJECT_GROUP_LIMIT_REACHED", Message: "Project group member limit reached"}

	// Class Membership-related
	ErrMembershipNotFound     = AppError{Code: "MEMBERSHIP_NOT_FOUND", Message: "Membership not found"}
	ErrAlreadyMember          = AppError{Code: "ALREADY_MEMBER", Message: "User is already a member of this community"}
	ErrMembershipCreateFailed = AppError{Code: "MEMBERSHIP_CREATE_FAILED", Message: "Failed to create membership"}
	ErrMembershipDeleteFailed = AppError{Code: "MEMBERSHIP_DELETE_FAILED", Message: "Failed to delete membership"}
	ErrInvalidMembershipData  = AppError{Code: "INVALID_MEMBERSHIP_DATA", Message: "Invalid membership data"}

	// Messaging-related
	ErrChannelNotFound = AppError{Code: "CHANNEL_NOT_FOUND", Message: "Channel not found"}
	ErrNoMessageFound  = AppError{Code: "NO_MESSAGE_FOUND", Message: "No message found"}

	// Classroom-related
	ErrClassroomNotFound = AppError{Code: "CLASSROOM_NOT_FOUND", Message: "Classroom not found"}
	ErrRoundNotFound     = AppError{Code: "ROUND_NOT_FOUND", Message: "Round not found"}
	ErrPostNotFound      = AppError{Code: "POST_NOT_FOUND", Message: "Post not found"}

	// Invitation-related
	ErrInvitationNotFound         = AppError{Code: "INVITATION_NOT_FOUND", Message: "Invitation not found"}
	ErrInvitationExpired          = AppError{Code: "INVITATION_EXPIRED", Message: "Invitation has expired"}
	ErrInvitationAlreadyExists    = AppError{Code: "INVITATION_ALREADY_EXISTS", Message: "Invitation already exists for this email"}
	ErrInvitationAlreadyProcessed = AppError{Code: "INVITATION_ALREADY_PROCESSED", Message: "Invitation has already been accepted or rejected"}
	ErrAlreadyInClassroom         = AppError{Code: "ALREADY_IN_CLASSROOM", Message: "User is already a member of this classroom"}
	ErrClassroomFull              = AppError{Code: "CLASSROOM_FULL", Message: "Classroom has reached maximum capacity"}
	ErrEmailMismatch              = AppError{Code: "EMAIL_MISMATCH", Message: "Email does not match invitation"}

	// Classroom Join-related
	ErrCodeInvalid                 = AppError{Code: "CODE_INVALID", Message: "Invalid invitation code"}
	ErrStudentCodeNotInWhitelist   = AppError{Code: "STUDENT_CODE_NOT_IN_WHITELIST", Message: "Student code not in whitelist"}
	ErrStudentCodeAlreadyUsed      = AppError{Code: "STUDENT_CODE_ALREADY_USED", Message: "Student code already used in this classroom"}
	ErrInvalidEmailDomain          = AppError{Code: "INVALID_EMAIL_DOMAIN", Message: "Email domain not allowed for this classroom"}
	ErrJoinRequestNotFound         = AppError{Code: "JOIN_REQUEST_NOT_FOUND", Message: "Join request not found"}
	ErrJoinRequestAlreadyProcessed = AppError{Code: "JOIN_REQUEST_ALREADY_PROCESSED", Message: "Join request has already been processed"}
	ErrAlreadyCoLecturer           = AppError{Code: "ALREADY_CO_LECTURER", Message: "User is already a co-lecturer of this classroom"}
	ErrAlreadyLecturer             = AppError{Code: "ALREADY_LECTURER", Message: "User is already the lecturer of this classroom"}
	ErrCannotInviteSelf            = AppError{Code: "CANNOT_INVITE_SELF", Message: "Cannot invite yourself"}
)
