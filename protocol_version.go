package cas

// ProtocolVersion represents a CAS protocol version
type ProtocolVersion uint8

const (
	// CASVersionUndefined represents the CAS protocol version to use is not set
	CASVersionUndefined ProtocolVersion = iota

	// CASVersion1 represents using the CAS 1.0 protocol
	CASVersion1

	// CASVersion2 represents using the CAS 2.0 protocol
	CASVersion2

	// CASVersion3 represents using the CAS 3.0 protocol
	CASVersion3
)
