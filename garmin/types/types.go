package types

import (
	"fmt"
	"strings"
)

type TicketRequest struct {
	Host string `json:"host"`
	// -- request params
	ClientId string `json:"clientId"`
	Locale   string `json:"locale"`
	Service  string `json:"service"`
}

type ActionTokenResponse struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	ExpiresIn             int64  `json:"expires_in"`
	Scope                 string `json:"scope"`
	Jti                   string `json:"jti"`
	RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
}

// return https://sso.garmin.cn/portal/api/login?clientId=GarminConnect&locale=zh-CN&service=https%3A%2F%2Fconnect.garmin.cn%2Fmodern
func (t *TicketRequest) Url(host string) string {
	return fmt.Sprintf(
		"https://sso.%s/portal/api/login?clientId=%s&locale=%s&service=%s",
		host,
		t.ClientId,
		t.Locale,
		t.Service,
	)
}

func NewTicketRequest(host string) *TicketRequest {
	return &TicketRequest{
		Host:     host,
		ClientId: "GarminConnect",
		Locale:   "zh-CN",
		Service:  fmt.Sprintf("https://connect.%s/modern", host),
	}
}

type ActionTicketResponse struct {
	ServiceURL           string         `json:"serviceURL"`
	ServiceTicketID      string         `json:"serviceTicketId"`
	ResponseStatus       ResponseStatus `json:"responseStatus"`
	CustomerMfaInfo      interface{}    `json:"customerMfaInfo"`
	ConsentTypeList      interface{}    `json:"consentTypeList"`
	CAPTCHAAlreadyPassed bool           `json:"captchaAlreadyPassed"`
}

type ResponseStatus struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	HTTPStatus string `json:"httpStatus"`
}

func (t *ActionTicketResponse) HomeUrl() string {
	return fmt.Sprintf("%s?ticket=%s", t.ServiceURL, t.ServiceTicketID)
}

type ActionUploadResponse struct {
	DetailedImportResult DetailedImportResult `json:"detailedImportResult"`
}

type DetailedImportResult struct {
	UploadID       any              `json:"uploadId"`
	UploadUUID     UploadUUID       `json:"uploadUuid"`
	Owner          int64            `json:"owner"`
	FileSize       int64            `json:"fileSize"`
	ProcessingTime int64            `json:"processingTime"`
	CreationDate   string           `json:"creationDate"`
	IPAddress      string           `json:"ipAddress"`
	FileName       string           `json:"fileName"`
	Report         interface{}      `json:"report"`
	Failures       []ActivityStatus `json:"failures"`
	Successes      []ActivityStatus `json:"successes"`
}

type UploadUUID struct {
	Uuid string `json:"uuid"`
}

type ActivityStatus struct {
	InternalID int64     `json:"internalId"`
	ExternalID string    `json:"externalId"`
	Messages   []Message `json:"messages"`
}

type Message struct {
	Code    int64  `json:"code"`
	Content string `json:"content"`
}

func (aur *ActionUploadResponse) Failures() []ActivityStatus {
	return aur.DetailedImportResult.Failures
}

func (aur *ActionUploadResponse) Fails() bool {
	return len(aur.Failures()) > 0
}

func (aur *ActionUploadResponse) Successes() []ActivityStatus {
	return aur.DetailedImportResult.Successes
}

func (aur *ActionUploadResponse) Success() bool {
	uploadId := aur.DetailedImportResult.UploadID
	uploadIdOk := false
	if uploadIdNum := uploadId.(float64); uploadIdNum > 0 {
		uploadIdOk = true
	}
	if uploadIdStr := uploadId.(string); uploadIdStr != "" {
		uploadIdOk = true
	}

	return ((uploadIdOk) &&
		!strings.EqualFold(
			aur.DetailedImportResult.UploadUUID.Uuid,
			"",
		)) || len(aur.DetailedImportResult.Successes) > 0
}
