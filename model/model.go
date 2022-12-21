package model

type KafkaModel struct {
	RefId   *string `json:"refId,omitempty"`
	State   *string `json:"state,omitempty"`
	Message *string `json:"message,omitempty"`
}

type Header struct {
	AppId       string `json:"appId"`
	AccessToken string `json:"accessToken"`
	ReferenceId string `json:"referenceId"`
	ServiceId   string `json:"serviceId"`
}
