package myutil

//ResoultReadFile :
type ResoultReadFile struct {
	FileName string
	Context  string
}

//Message : use this message struct to send command between nets
type Message struct {
	Mothed string      `json:"Mothed"`
	Data   interface{} `json:"Data"`
}

//OutOSArgs :
type OutOSArgs struct {
	HTTPPort string
	NetProt  string
	LogPath  string
}
