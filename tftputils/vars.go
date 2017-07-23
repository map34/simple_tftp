package tftputils

//  opcode  operation
// 	1     Read request (RRQ)
// 	2     Write request (WRQ)
// 	3     Data (DATA)
// 	4     Acknowledgment (ACK)
// 	5     Error (ERROR)
const (
	UNKNOWNOP = iota
	RRQ
	WRQ
	DATA
	ACK
	ERROR
)

const (
	UnknownErr = iota
	FileNotFoundErr
	AccessViolationErr
	DiskFullErr
	IllegalOpErr
	UnknownTransferIDErr
	FileExistsErr
	NoSuchUserErr
)

const (
	SmallestBlockSize = 512
	Uint64BytesNum  = 8
)

const (
	OCTET = "octet"
)
