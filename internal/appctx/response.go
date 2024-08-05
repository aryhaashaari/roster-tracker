// Package appctx
package appctx

import (
	"encoding/json"
	"fmt"
	"sync"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/msgx"
)

var (
	rsp    *Response
	oneRsp sync.Once
)

// Response presentation contract object
type Response struct {
	Code      int         `json:"code"`
	Status    string      `json:"status,omitempty"`
	Entity    string      `json:"entity,omitempty"`
	State     string      `json:"state,omitempty"`
	Message   interface{} `json:"message,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	lang      string      `json:"-"`
	Meta      interface{} `json:"meta,omitempty"`
	msgKey    string
	requestId string `json:"-"`
}

// MetaData represent meta data response for multi data
type MetaData struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalPage  int64 `json:"total_page"`
	TotalCount int64 `json:"total_count"`
}

// GetMessage method to transform response name var to message detail
func (r *Response) GetMessage() string {
	return msgx.Get(r.msgKey, r.lang).Text()
}

// Generate setter message
func (r *Response) Generate() *Response {
	if r.lang == "" {
		r.lang = consts.LangDefault
	}
	msg := msgx.Get(r.msgKey, r.lang)
	if r.Message == nil {
		r.Message = msg.Text()
	}

	if r.Code == 0 {
		r.Code = msg.Status()
	}

	return r
}

// WithStatus setter status response
func (r *Response) WithStatus(s string) *Response {
	r.Status = s
	r.State = fmt.Sprintf("%s%s", r.Entity, s)
	return r
}

// WithEntity setter entity response
func (r *Response) WithEntity(e string) *Response {
	r.Entity = e
	return r.WithStatus("SUCCESS")
}

// WithState setter state response
func (r *Response) WithState(s string) *Response {
	r.State = s
	return r
}

// WithCode setter response var name
func (r *Response) WithCode(c int) *Response {
	r.Code = c
	return r
}

// WithData setter data response
func (r *Response) WithData(v any) *Response {
	r.Data = v
	return r
}

// WithError setter error messages
func (r *Response) WithError(v any) *Response {
	r.Errors = v
	return r
}

func (r *Response) WithMsgKey(v string) *Response {
	r.msgKey = v
	return r
}

// WithMeta setter meta data response
func (r *Response) WithMeta(v any) *Response {
	r.Meta = v
	return r
}

// WithLang setter language response
func (r *Response) WithLang(v string) *Response {
	r.lang = v
	return r
}

// WithMessage setter custom message response
func (r *Response) WithMessage(v any) *Response {
	if v != nil {
		r.Message = v
	}

	return r
}

func (r *Response) Byte() []byte {
	if r.Code == 0 || r.Message == nil {
		r.Generate()
	}

	b, _ := json.Marshal(r)
	return b
}

func (r *Response) WithRequestID(v string) *Response {
	r.requestId = v
	return r
}

func (r *Response) RequestID() string {
	return r.requestId
}

// NewResponse initialize response
func NewResponse() *Response {
	oneRsp.Do(func() {
		rsp = &Response{
			Status: "SUCCESS",
		}
	})

	// clone response
	x := *rsp

	return &x
}
