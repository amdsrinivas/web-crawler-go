package models

type Resource struct {
	ResourceAddress string
	ResourceType    string
}

type ResourceData struct {
	ResourceAddress string
	Data            []byte
}
