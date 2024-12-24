package models

type HeartBeat struct {
	Replication int      `json:"Replication"`
	MainAddres  string   `json:"MainAddres"`
	Address     []string `json:"ServerAddress"`
}
type GetReq struct {
	UUID4 string `json:"UUID4"`
}
type DelReq struct {
	UUID4 string `json:"UUID4"`
}
type SetReq struct {
	UUID4 string `json:"UUID4"`
	Value string `json:"Value"`
}

type Flags struct {
	H,
	P,
	HBD,
	PBD *string
	R         *int
	Address   string
	BDaddress string
}
