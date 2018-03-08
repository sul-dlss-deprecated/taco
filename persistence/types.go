package persistence

import "github.com/go-openapi/strfmt"

// AtContext is our primary context
const AtContext = "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld"

var (
	// FileType is the type for file resources
	FileType = strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-file.jsonld")
	// FilesetType is the type for fileset resources
	FilesetType = strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-fileset.jsonld")
	// SequenceType is the type for sequence resources
	SequenceType = strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-sequence.jsonld")
	// PrimarySequenceType is the type for primary sequence resources
	PrimarySequenceType = strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-primary-sequence.jsonld")
)

// IsObjectOrCollection returns true if the object is a file or a collection
func IsObjectOrCollection(uri strfmt.URI) bool {
	return !(uri == FileType ||
		uri == FilesetType ||
		uri == SequenceType ||
		uri == PrimarySequenceType)
}
