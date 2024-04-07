package models

type CommonApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
}

func CreateCommonSuccessResponse(Data interface{}) CommonApiResponse {
	return CommonApiResponse{Success: true, Data: Data, Error: ""}
}

func CreateCommonErrorResponse(Error string) CommonApiResponse {
	return CommonApiResponse{Success: false, Data: nil, Error: Error}
}
