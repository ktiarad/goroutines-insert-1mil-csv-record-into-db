package model

type Domain struct {
	GlobalRank     int
	TldRank        int
	Domain         string
	TLD            string
	RefSubNets     int
	RefIPs         int
	IDN_Domain     string
	IDN_TLD        string
	PrevGlobalRank int
	PrevTldRank    int
	PrevRefSubNets int
	PrevRefIPs     int
}
