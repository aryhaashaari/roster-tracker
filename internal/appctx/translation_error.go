package appctx

import (
	"context"
	"database/sql"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/provider/providererror"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/logger"
	"net/http"
)

type ResponderTranslationErrorExtender func(ctx context.Context, err error, fl ...logger.Field) *Response

type ResponderTranslationErrorOpt struct {
	loggerFields                      []logger.Field
	responderTranslationErrorExtender ResponderTranslationErrorExtender
	translateDBError                  bool
}

type ResponderOpt func(opt *ResponderTranslationErrorOpt)

func (r *Response) WithLoggerField(lf ...logger.Field) ResponderOpt {
	return func(opt *ResponderTranslationErrorOpt) {
		opt.loggerFields = lf
	}
}

func (r *Response) WithExtender(extender ResponderTranslationErrorExtender) ResponderOpt {
	return func(opt *ResponderTranslationErrorOpt) {
		opt.responderTranslationErrorExtender = extender
	}
}

func (r *Response) WithTranslateDBError() ResponderOpt {
	return func(opt *ResponderTranslationErrorOpt) {
		opt.translateDBError = true
	}
}

func (r *Response) ResponderTranslationError(ctx context.Context, err error, opts ...ResponderOpt) *Response {
	opt := ResponderTranslationErrorOpt{
		loggerFields:                      make([]logger.Field, 0),
		responderTranslationErrorExtender: nil,
		translateDBError:                  false,
	}

	for _, o := range opts {
		o(&opt)
	}

	fl := opt.loggerFields

	r.WithStatus("ERROR")

	if extender := opt.responderTranslationErrorExtender; extender != nil {
		res := extender(ctx, err, fl...)
		if res != nil {
			return res
		}
	}

	errCause := errors.Cause(err)
	switch errCause {

	case consts.ErrInvalidJSON:
		logger.WarnWithContext(ctx, err, fl...)
		return r.WithCode(http.StatusBadRequest).WithMessage("Invalid JSON Request")

	case sql.ErrNoRows:
		logger.WarnWithContext(ctx, err, fl...)
		return r.WithCode(http.StatusNotFound).WithMessage("Data not found")

	default:
		switch errCauser := errCause.(type) {

		case validation.Errors:
			logger.WarnWithContext(ctx, err, fl...)
			return r.WithCode(http.StatusUnprocessableEntity).WithError(errCause).WithMessage("Validation Error(s)")

		case consts.Error:
			logger.WarnWithContext(ctx, err, fl...)
			return r.WithCode(http.StatusBadRequest).WithMessage(errCause.Error())

		case consts.ErrNotFound:
			logger.WarnWithContext(ctx, err, fl...)
			return r.WithCode(http.StatusNotFound).WithMessage(errCause.Error())

		case consts.ErrAlreadyExist:
			logger.WarnWithContext(ctx, err, fl...)
			return r.WithCode(http.StatusConflict).WithMessage(errCause.Error())

		case providererror.Error:
			logger.WarnWithContext(ctx, err, fl...)
			return r.WithCode(errCauser.Code).WithMessage(errCauser.Message).WithError(errCauser.Errors)
		}
	}

	if opt.translateDBError {
		res := r.responderTranslationDB(ctx, err, fl...)
		if res != nil {
			return res
		}
	}

	logger.ErrorWithContext(ctx, err, fl...)
	return r.WithCode(http.StatusInternalServerError).WithMessage(http.StatusText(http.StatusInternalServerError))
}

func (r *Response) responderTranslationDB(ctx context.Context, err error, fl ...logger.Field) *Response {
	switch errt := errors.Cause(err).(type) {
	case *pq.Error:
		switch errt.Code {
		case "23505":
			logger.WarnWithContext(ctx, err, fl...)
			return r.WithCode(http.StatusBadRequest).WithMessage(errt.Detail)
		}
	}

	return nil
}
