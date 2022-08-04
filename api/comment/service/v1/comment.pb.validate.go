// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: comment/service/v1/comment.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// define the regex for a UUID once up-front
var _comment_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on CreateCommentDraftReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateCommentDraftReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateCommentDraftReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateCommentDraftReqMultiError, or nil if none found.
func (m *CreateCommentDraftReq) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateCommentDraftReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetUuid()); err != nil {
		err = CreateCommentDraftReqValidationError{
			field:  "Uuid",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CreateCommentDraftReqMultiError(errors)
	}

	return nil
}

func (m *CreateCommentDraftReq) _validateUuid(uuid string) error {
	if matched := _comment_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// CreateCommentDraftReqMultiError is an error wrapping multiple validation
// errors returned by CreateCommentDraftReq.ValidateAll() if the designated
// constraints aren't met.
type CreateCommentDraftReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateCommentDraftReqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateCommentDraftReqMultiError) AllErrors() []error { return m }

// CreateCommentDraftReqValidationError is the validation error returned by
// CreateCommentDraftReq.Validate if the designated constraints aren't met.
type CreateCommentDraftReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateCommentDraftReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateCommentDraftReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateCommentDraftReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateCommentDraftReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateCommentDraftReqValidationError) ErrorName() string {
	return "CreateCommentDraftReqValidationError"
}

// Error satisfies the builtin error interface
func (e CreateCommentDraftReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateCommentDraftReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateCommentDraftReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateCommentDraftReqValidationError{}

// Validate checks the field values on CreateCommentDraftReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateCommentDraftReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateCommentDraftReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateCommentDraftReplyMultiError, or nil if none found.
func (m *CreateCommentDraftReply) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateCommentDraftReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return CreateCommentDraftReplyMultiError(errors)
	}

	return nil
}

// CreateCommentDraftReplyMultiError is an error wrapping multiple validation
// errors returned by CreateCommentDraftReply.ValidateAll() if the designated
// constraints aren't met.
type CreateCommentDraftReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateCommentDraftReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateCommentDraftReplyMultiError) AllErrors() []error { return m }

// CreateCommentDraftReplyValidationError is the validation error returned by
// CreateCommentDraftReply.Validate if the designated constraints aren't met.
type CreateCommentDraftReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateCommentDraftReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateCommentDraftReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateCommentDraftReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateCommentDraftReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateCommentDraftReplyValidationError) ErrorName() string {
	return "CreateCommentDraftReplyValidationError"
}

// Error satisfies the builtin error interface
func (e CreateCommentDraftReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateCommentDraftReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateCommentDraftReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateCommentDraftReplyValidationError{}

// Validate checks the field values on GetLastCommentDraftReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetLastCommentDraftReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetLastCommentDraftReq with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetLastCommentDraftReqMultiError, or nil if none found.
func (m *GetLastCommentDraftReq) ValidateAll() error {
	return m.validate(true)
}

func (m *GetLastCommentDraftReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetUuid()); err != nil {
		err = GetLastCommentDraftReqValidationError{
			field:  "Uuid",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetLastCommentDraftReqMultiError(errors)
	}

	return nil
}

func (m *GetLastCommentDraftReq) _validateUuid(uuid string) error {
	if matched := _comment_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// GetLastCommentDraftReqMultiError is an error wrapping multiple validation
// errors returned by GetLastCommentDraftReq.ValidateAll() if the designated
// constraints aren't met.
type GetLastCommentDraftReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetLastCommentDraftReqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetLastCommentDraftReqMultiError) AllErrors() []error { return m }

// GetLastCommentDraftReqValidationError is the validation error returned by
// GetLastCommentDraftReq.Validate if the designated constraints aren't met.
type GetLastCommentDraftReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetLastCommentDraftReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetLastCommentDraftReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetLastCommentDraftReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetLastCommentDraftReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetLastCommentDraftReqValidationError) ErrorName() string {
	return "GetLastCommentDraftReqValidationError"
}

// Error satisfies the builtin error interface
func (e GetLastCommentDraftReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetLastCommentDraftReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetLastCommentDraftReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetLastCommentDraftReqValidationError{}

// Validate checks the field values on GetLastCommentDraftReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetLastCommentDraftReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetLastCommentDraftReply with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetLastCommentDraftReplyMultiError, or nil if none found.
func (m *GetLastCommentDraftReply) ValidateAll() error {
	return m.validate(true)
}

func (m *GetLastCommentDraftReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if len(errors) > 0 {
		return GetLastCommentDraftReplyMultiError(errors)
	}

	return nil
}

// GetLastCommentDraftReplyMultiError is an error wrapping multiple validation
// errors returned by GetLastCommentDraftReply.ValidateAll() if the designated
// constraints aren't met.
type GetLastCommentDraftReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetLastCommentDraftReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetLastCommentDraftReplyMultiError) AllErrors() []error { return m }

// GetLastCommentDraftReplyValidationError is the validation error returned by
// GetLastCommentDraftReply.Validate if the designated constraints aren't met.
type GetLastCommentDraftReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetLastCommentDraftReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetLastCommentDraftReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetLastCommentDraftReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetLastCommentDraftReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetLastCommentDraftReplyValidationError) ErrorName() string {
	return "GetLastCommentDraftReplyValidationError"
}

// Error satisfies the builtin error interface
func (e GetLastCommentDraftReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetLastCommentDraftReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetLastCommentDraftReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetLastCommentDraftReplyValidationError{}

// Validate checks the field values on SendCommentReq with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *SendCommentReq) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SendCommentReq with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in SendCommentReqMultiError,
// or nil if none found.
func (m *SendCommentReq) ValidateAll() error {
	return m.validate(true)
}

func (m *SendCommentReq) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if err := m._validateUuid(m.GetUuid()); err != nil {
		err = SendCommentReqValidationError{
			field:  "Uuid",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if ip := net.ParseIP(m.GetIp()); ip == nil {
		err := SendCommentReqValidationError{
			field:  "Ip",
			reason: "value must be a valid IP address",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return SendCommentReqMultiError(errors)
	}

	return nil
}

func (m *SendCommentReq) _validateUuid(uuid string) error {
	if matched := _comment_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// SendCommentReqMultiError is an error wrapping multiple validation errors
// returned by SendCommentReq.ValidateAll() if the designated constraints
// aren't met.
type SendCommentReqMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SendCommentReqMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SendCommentReqMultiError) AllErrors() []error { return m }

// SendCommentReqValidationError is the validation error returned by
// SendCommentReq.Validate if the designated constraints aren't met.
type SendCommentReqValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SendCommentReqValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SendCommentReqValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SendCommentReqValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SendCommentReqValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SendCommentReqValidationError) ErrorName() string { return "SendCommentReqValidationError" }

// Error satisfies the builtin error interface
func (e SendCommentReqValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSendCommentReq.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SendCommentReqValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SendCommentReqValidationError{}