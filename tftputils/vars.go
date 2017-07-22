package tftputils

//  opcode  operation
// 	1     Read request (RRQ)
// 	2     Write request (WRQ)
// 	3     Data (DATA)
// 	4     Acknowledgment (ACK)
// 	5     Error (ERROR)
const (
	UNKNOWNOP = 0
	RRQ       = 1
	WRQ       = 2
	DATA      = 3
	ACK       = 4
	ERROR     = 5
)

const (
	SMALLESTBLOCK  = 512
	UINT64BYTESNUM = 8
)

const (
	OCTET = "octet"
)
