package config

import (
	"ex00/bd/internal/models"
	"flag"
)

var Flager models.Flags
var HB models.HeartBeat = models.HeartBeat{}

func InitFlags() {
	Flager.H = flag.String("H", "127.0.0.1", "Print host")
	Flager.P = flag.String("P", "8765", "Print port")
	Flager.HBD = flag.String("HBD", "", "Print port")
	Flager.PBD = flag.String("PBD", "", "Print port")
	Flager.R = flag.Int("R", 2, "Print port")

	flag.Parse()
	Flager.Address = *Flager.H + ":" + *Flager.P
	HB.MainAddres = *Flager.H + ":" + *Flager.P
	Flager.BDaddress = *Flager.HBD + ":" + *Flager.PBD
	HB.Address = append(HB.Address, Flager.Address)
	HB.Replication = *Flager.R

}

func CheckConfig() bool {
	return Flager.BDaddress != ":"
}
