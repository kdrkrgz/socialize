package utils

type ApiResponse struct {
	IsSuccess bool              `default:"true" json:"isSuccess"`
	Data      interface{}       `json:"data"`
	Errors    map[string]string `json:"errors"`
}

func SuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		IsSuccess: true,
		Data:      data,
		Errors:    nil,
	}
}

func ErrorResponse(errs map[string]string) ApiResponse {
	return ApiResponse{
		IsSuccess: false,
		Data:      nil,
		Errors:    errs,
	}
}

func ApiErrors(key string) map[string]string {
	var apiErrs []map[string]string
	apiErrs = append(apiErrs, map[string]string{"ServiceError": "Something went wrong"})
	apiErrs = append(apiErrs, authErrs()...)
	return func(key string) map[string]string {
		for i := range apiErrs {
			if val, ok := apiErrs[i][key]; ok {
				return map[string]string{key: val}
			} else {
				continue
			}
		}
		return apiErrs[0]
	}(key)
}

func authErrs() []map[string]string {
	var errs []map[string]string
	errs = append(errs, map[string]string{"UnAuthorized": "You are not logged in"})
	errs = append(errs, map[string]string{"InvalidToken": "Your token is invalid"})
	errs = append(errs, map[string]string{"Forbidden": "You are not authorized to access this resource"})
	errs = append(errs, map[string]string{"UserNotExist": "The user belonging to this token no longer exists"})
	return errs
}
