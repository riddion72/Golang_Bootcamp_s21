package models

type HeartBeat struct {
	Replication int      `json:"Replication"`
	MainAddres  string   `json:"MainAddres"`
	Address     []string `json:"ServerAddress"`
}

type GetRequestModel struct {
	UUID4 string `json:"UUID4"`
}

type SetRequestModel struct {
	UUID4 string `json:"UUID4"`
	Value string `json:"Value"`
}

type Flags struct {
	H,
	P *string
	Addres string
}
