package model

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	CodeInvalidParameter = "INVALID_PARAMETER"
	CodeNotFound         = "NOT_FOUND"
	CodeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	CodeInternalError    = "INTERNAL_ERROR"
	CodeUpstreamError    = "UPSTREAM_ERROR"
	CodeUpstreamTimeout  = "UPSTREAM_TIMEOUT"
)
