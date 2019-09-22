package response

import (
	"encoding/json"

	HTTP_CODE "bitbucket.org/sparkmaker/gohelper/http/code"
)

const (
	// Internal Code
	ICBadRequest          = 10
	ICUnauthorized        = 20
	ICNotFound            = 30
	ICGone                = 30
	ICInternalServerError = 40
	ICServiceUnavailable  = 50
)

type ResponseWrapper struct {
	Content []byte
	Code    int
}

func DataResponse(data interface{}, httpCode int) ResponseWrapper {
	t := map[string]interface{}{
		"data": data,
		"code": httpCode,
	}
	arrByte, err := json.Marshal(&t)
	if err != nil {
		return ErrorResponse("Internal parsing failure.", []string{}, ICInternalServerError, HTTP_CODE.InternalServerError)
	}
	return ResponseWrapper{arrByte, httpCode}
}

func ErrorResponse(message string, errorInfo []string, internalCode int, httpCode int) ResponseWrapper {
	t := map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"code":    internalCode,
			"info":    errorInfo,
		},
		"code": httpCode,
	}
	arrByte, err := json.Marshal(&t)
	if err != nil {
		return ErrorResponse("Internal parsing failure.", []string{}, ICInternalServerError, HTTP_CODE.InternalServerError)
	}
	return ResponseWrapper{arrByte, httpCode}
}

func CustomErrorResponse(errorData interface{}, httpCode int) ResponseWrapper {
	t := map[string]interface{}{
		"error": errorData,
		"code":  httpCode,
	}
	arrByte, err := json.Marshal(&t)
	if err != nil {
		return ErrorResponse("Internal parsing failure.", []string{}, ICInternalServerError, HTTP_CODE.InternalServerError)
	}
	return ResponseWrapper{arrByte, httpCode}
}

func OK(data interface{}) ResponseWrapper {
	return DataResponse(data, HTTP_CODE.OK)
}

func CREATED(data interface{}) ResponseWrapper {
	return DataResponse(data, HTTP_CODE.Created)
}

func BAD_REQUEST(info []string) ResponseWrapper {
	return ErrorResponse("Bad Request.", info, ICBadRequest, HTTP_CODE.BadRequest)
}

func UNAUTHORIZED(info []string) ResponseWrapper {
	return ErrorResponse("Unauthorized.", info, ICUnauthorized, HTTP_CODE.Unauthorized)
}

func NOT_FOUND(info []string) ResponseWrapper {
	return ErrorResponse("Not Found.", info, ICNotFound, HTTP_CODE.NotFound)
}

func INTERNAL_SERVER_ERROR(info []string) ResponseWrapper {
	return ErrorResponse("Internal Server Error.", info, ICInternalServerError, HTTP_CODE.InternalServerError)
}

func SERVICE_UNAVAILABLE(info []string) ResponseWrapper {
	return ErrorResponse("Service unavailable.", info, ICServiceUnavailable, HTTP_CODE.ServiceUnavailable)
}

func GONE(info []string) ResponseWrapper {
	return ErrorResponse("Gone.", info, ICGone, HTTP_CODE.Gone)
}
