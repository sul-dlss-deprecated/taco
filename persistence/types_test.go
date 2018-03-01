package persistence

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

func TestIsObjectOrCollection(t *testing.T) {
	assert.True(t, IsObjectOrCollection(strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-exhibit.jsonld")))
	assert.True(t, IsObjectOrCollection(strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-photograph.jsonld")))
	assert.False(t, IsObjectOrCollection(strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-sequence.jsonld")))
	assert.False(t, IsObjectOrCollection(strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-primary-sequence.jsonld")))
	assert.False(t, IsObjectOrCollection(strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-fileset.jsonld")))
	assert.False(t, IsObjectOrCollection(strfmt.URI("http://sdr.sul.stanford.edu/models/sdr3-file.jsonld")))
}
