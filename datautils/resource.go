package datautils

import (
	"bytes"
	"fmt"
)

// FileType is the type URI for files
const FileType = "http://sdr.sul.stanford.edu/models/sdr3-file.jsonld"

// FilesetType is the type URI for filesets
const FilesetType = "http://sdr.sul.stanford.edu/models/sdr3-fileset.jsonld"

// ObjectTypes is the list of object subtype URIs
var ObjectTypes = []string{
	"http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-3d.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-agreement.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-book.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-document.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-geo.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-image.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-page.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-photograph.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-manuscript.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-map.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-media.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-track.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-webarchive-binary.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-webarchive-seed.jsonld",
}

// CollectionTypes is the list of object subtype URIs
var CollectionTypes = []string{
	"http://sdr.sul.stanford.edu/models/sdr3-collection.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-curated-collection.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-user-collection.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-exhibit.jsonld",
	"http://sdr.sul.stanford.edu/models/sdr3-series.jsonld",
}

// Resource represents the resource as it exists in the persistence layer
// this is very similar to models.Resource, but COULD vary, so we should
// keep them separated
type Resource struct {
	JSON JSONObject
}

// NewResource creates a new resource instance
func NewResource(data map[string]interface{}) *Resource {
	if data == nil {
		data = map[string]interface{}{}
	}
	return &Resource{JSON: JSONObject(data)}
}

// ID returns the document's identifier
func (d *Resource) ID() string {
	return d.JSON.GetS("tacoIdentifier")
}

// ExternalIdentifier returns the document's external identifier (DRUID or UUID)
func (d *Resource) ExternalIdentifier() string {
	return d.JSON.GetS("externalIdentifier")
}

// Version returns the document's version
func (d *Resource) Version() int {
	return int(d.JSON.GetF("version"))
}

// Type returns the document's type
func (d *Resource) Type() string {
	return d.JSON.GetS("@type")
}

// MimeType returns the document's MIME type
func (d *Resource) MimeType() string {
	return d.JSON.GetS("hasMimeType")
}

// WithFileLocation sets the location of the binary.
func (d *Resource) FileLocation() string {
	return d.JSON.GetS("file-location")
}

// Label returns the document's Label
func (d *Resource) Label() string {
	return d.JSON.GetS("label")
}

// IsFile returns true if the resource has the file type assertion
func (d *Resource) IsFile() bool {
	return d.Type() == FileType
}

// IsFileset returns true if the resource has the fileset type assertion
func (d *Resource) IsFileset() bool {
	return d.Type() == FilesetType
}

// IsObject returns true if the resource has an object type assertion
func (d *Resource) IsObject() bool {
	return contains(ObjectTypes, d.Type())
}

// IsCollection returns true if the resource has an object type assertion
func (d *Resource) IsCollection() bool {
	return contains(CollectionTypes, d.Type())
}

// WithID sets the document's primary key
func (d *Resource) WithID(id string) *Resource {
	d.JSON["tacoIdentifier"] = id
	return d
}

// WithType sets the document's type
func (d *Resource) WithType(atType string) *Resource {
	d.JSON["@type"] = atType
	return d
}

// WithExternalIdentifier sets the document's external identifier (DRUID or UUID)
func (d *Resource) WithExternalIdentifier(id string) *Resource {
	d.JSON["externalIdentifier"] = id
	return d
}

// WithMimeType sets the mime type. This should only be used on File resources
func (d *Resource) WithMimeType(mimeType string) *Resource {
	d.JSON["hasMimeType"] = mimeType
	return d
}

// WithFileLocation sets the location of the binary. This should only be used on File resources
func (d *Resource) WithFileLocation(fileLocation string) *Resource {
	d.JSON["file-location"] = fileLocation
	return d
}

// WithLabel sets the label.
func (d *Resource) WithLabel(label string) *Resource {
	d.JSON["label"] = label
	return d
}

// WithCurrentVersion sets the currentVersion flag
func (d *Resource) WithCurrentVersion(flag bool) *Resource {
	d.JSON["currentVersion"] = flag
	return d
}

// WithVersion sets the version
func (d *Resource) WithVersion(version int) *Resource {
	d.JSON["version"] = float64(version)
	return d
}

// WithPrecedingVersion sets the precedingVersion to
// the id passed (of the old version)
func (d *Resource) WithPrecedingVersion(id string) *Resource {
	d.JSON["precedingVersion"] = id
	return d
}

// WithFollowingVersion sets the followingVersion to
// the id passed (of the new version)
func (d *Resource) WithFollowingVersion(id string) *Resource {
	d.JSON["followingVersion"] = id
	return d
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Structural returns structural.isContainedBy
func (d *Resource) Structural() *JSONObject {
	return d.JSON.GetObj("structural")
}

// Identification returns the identification subschema
func (d *Resource) Identification() *JSONObject {
	return d.JSON.GetObj("identification")
}

// Administrative returns the administrative subschema
func (d *Resource) Administrative() *JSONObject {
	return d.JSON.GetObj("administrative")
}

func (d *Resource) String() string {
	buf := bytes.NewBufferString("<Resource")
	if d.JSON.HasKey("tacoIdentifier") {
		buf.WriteString(fmt.Sprintf(" id:'%s'", d.ID()))
	}
	if d.JSON.HasKey("@type") {
		buf.WriteString(fmt.Sprintf(" @type:'%s'", d.Type()))
	}
	buf.WriteString(">")
	return buf.String()
}
