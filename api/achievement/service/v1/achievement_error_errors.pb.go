// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package v1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsUnknownError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_UNKNOWN_ERROR.String() && e.Code == 500
}

func ErrorUnknownError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_UNKNOWN_ERROR.String(), fmt.Sprintf(format, args...))
}

func IsSetAchievementAgreeFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_SET_ACHIEVEMENT_AGREE_FAILED.String() && e.Code == 500
}

func ErrorSetAchievementAgreeFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_SET_ACHIEVEMENT_AGREE_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsSetAchievementViewFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_SET_ACHIEVEMENT_VIEW_FAILED.String() && e.Code == 500
}

func ErrorSetAchievementViewFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_SET_ACHIEVEMENT_VIEW_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsSetAchievementCollectFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_SET_ACHIEVEMENT_COLLECT_FAILED.String() && e.Code == 500
}

func ErrorSetAchievementCollectFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_SET_ACHIEVEMENT_COLLECT_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsSetAchievementFollowFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_SET_ACHIEVEMENT_FOLLOW_FAILED.String() && e.Code == 500
}

func ErrorSetAchievementFollowFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_SET_ACHIEVEMENT_FOLLOW_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsCancelAchievementAgreeFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_CANCEL_ACHIEVEMENT_AGREE_FAILED.String() && e.Code == 500
}

func ErrorCancelAchievementAgreeFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_CANCEL_ACHIEVEMENT_AGREE_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsCancelAchievementCollectFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_CANCEL_ACHIEVEMENT_COLLECT_FAILED.String() && e.Code == 500
}

func ErrorCancelAchievementCollectFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_CANCEL_ACHIEVEMENT_COLLECT_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsCancelAchievementFollowFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_CANCEL_ACHIEVEMENT_FOLLOW_FAILED.String() && e.Code == 500
}

func ErrorCancelAchievementFollowFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_CANCEL_ACHIEVEMENT_FOLLOW_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsGetAchievementListFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_GET_ACHIEVEMENT_LIST_FAILED.String() && e.Code == 500
}

func ErrorGetAchievementListFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_GET_ACHIEVEMENT_LIST_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsGetAchievementFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_GET_ACHIEVEMENT_FAILED.String() && e.Code == 500
}

func ErrorGetAchievementFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_GET_ACHIEVEMENT_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsAddAchievementScoreFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_ADD_ACHIEVEMENT_SCORE_FAILED.String() && e.Code == 500
}

func ErrorAddAchievementScoreFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_ADD_ACHIEVEMENT_SCORE_FAILED.String(), fmt.Sprintf(format, args...))
}

func IsReduceAchievementScoreFailed(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == AchievementErrorReason_REDUCE_ACHIEVEMENT_SCORE_FAILED.String() && e.Code == 500
}

func ErrorReduceAchievementScoreFailed(format string, args ...interface{}) *errors.Error {
	return errors.New(500, AchievementErrorReason_REDUCE_ACHIEVEMENT_SCORE_FAILED.String(), fmt.Sprintf(format, args...))
}
