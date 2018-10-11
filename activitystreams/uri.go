package activitystreams

import (
	"strings"
)

type (
	// IRI is a Internationalized Resource Identifiers (IRIs) RFC3987
	IRI URI

	// URI is a Uniform Resource Identifier (URI) RFC3986
	URI string
)

// String returns the String value of the URI object
func (u URI) String() string {
	return string(u)
}

// String returns the String value of the IRI object
func (i IRI) String() string {
	return string(i)
}

// GetLink returns a copy of itself
func (u URI) GetLink() URI {
	return u
}

// GetLink returns a URI type coercion of the IRI object
func (i IRI) GetLink() URI {
	return URI(i)
}

// UnmarshalJSON
func (u *URI) UnmarshalJSON(s []byte) error {
	*u = URI(strings.Trim(string(s), "\""))
	return nil
}

// UnmarshalJSON
func (i *IRI) UnmarshalJSON(s []byte) error {
	*i = IRI(strings.Trim(string(s), "\""))
	return nil
}

// UnmarshalText
func (u URI) UnmarshalText(s []byte) error {
	u = URI(strings.Trim(string(s), "\""))
	return nil
}

// UnmarshalText
func (i IRI) UnmarshalText(s []byte) error {
	i = IRI(strings.Trim(string(s), "\""))
	return nil
}

// IsObject
func (u URI) GetID() *ObjectID {
	return nil
}

// IsObject
func (i IRI) GetID() *ObjectID {
	return nil
}

// GetType
func (u URI) GetType() ActivityVocabularyType {
	return LinkType
}

// GetType
func (i IRI) GetType() ActivityVocabularyType {
	return LinkType
}

// IsLink
func (u URI) IsLink() bool {
	return true
}

// IsLink
func (i IRI) IsLink() bool {
	return true
}

// IsObject
func (u URI) IsObject() bool {
	return false
}

// IsObject
func (i IRI) IsObject() bool {
	return false
}