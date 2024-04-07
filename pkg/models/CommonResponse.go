package models

type MetaData struct {
	TotalResults int `json:"totalResults"`
	PageSize     int `json:"pageSize"`
	Page         int `json:"page"`
}
type CommonApiResponse struct {
	Success  bool        `json:"success"`
	Data     interface{} `json:"data"`
	Error    string      `json:"error"`
	MetaData MetaData    `json:"metaData"`
}

func CreateCommonSuccessResponse(Data interface{}) CommonApiResponse {
	return CommonApiResponse{Success: true, Data: Data, Error: ""}
}

func CreateCommonSuccessWithMetaDataResponse(Data interface{}, metaData MetaData) CommonApiResponse {
	return CommonApiResponse{Success: true, Data: Data, Error: "", MetaData: metaData}
}

func CreateCommonErrorResponse(Error string) CommonApiResponse {
	return CommonApiResponse{Success: false, Data: nil, Error: Error}
}
