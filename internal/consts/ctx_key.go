// Package consts
package consts

const (
	// CtxUserCode const
	CtxUserCode = iota
	// CtxUserEmail const
	CtxUserEmail
	// CtxUserPhone const
	CtxUserPhone
	// CtxIP const
	CtxIP
	// CtxLang const
	CtxLang
	// CtxUserInfo const
	CtxUserInfo

	CtxtStartTime
)

const (
	ContextRequestID     = `request-id`
	ContextStartTime     = `start-time`
	ContextRequestIp     = `request-ip`
	ContextRequestMethod = `request-method`
	ContextRequestPath   = `request-path`
)
