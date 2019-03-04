package activitypub

import as "github.com/go-ap/activitystreams"

// Endpoints a json object which maps additional (typically server/domain-wide)
// endpoints which may be useful either for this actor or someone referencing this actor.
// This mapping may be nested inside the actor document as the value or may be a link to
// a JSON-LD document with these properties.
type Endpoints struct {
	// UploadMedia Upload endpoint URI for this user for binary data.
	UploadMedia as.Item `jsonld:"uploadMedia,omitempty"`
	// OauthAuthorizationEndpoint Endpoint URI so this actor's clients may access remote ActivityStreams objects which require authentication
	// to access. To use this endpoint, the client posts an x-www-form-urlencoded id parameter with the value being
	// the id of the requested ActivityStreams object.
	OauthAuthorizationEndpoint as.Item `jsonld:"oauthAuthorizationEndpoint,omitempty"`
	// OauthTokenEndpoint If OAuth 2.0 bearer tokens [RFC6749] [RFC6750] are being used for authenticating client to server interactions,
	// this endpoint specifies a URI at which a browser-authenticated user may obtain a new authorization grant.
	OauthTokenEndpoint as.Item `jsonld:"oauthTokenEndpoint,omitempty"`
	// ProvideClientKey  If OAuth 2.0 bearer tokens [RFC6749] [RFC6750] are being used for authenticating client to server interactions,
	// this endpoint specifies a URI at which a client may acquire an access token.
	ProvideClientKey as.Item `jsonld:"provideClientKey,omitempty"`
	// SignClientKey If Linked Data Signatures and HTTP Signatures are being used for authentication and authorization,
	// this endpoint specifies a URI at which browser-authenticated users may authorize a client's public
	// key for client to server interactions.
	SignClientKey as.Item `jsonld:"signClientKey,omitempty"`
	// SharedInbox If Linked Data Signatures and HTTP Signatures are being used for authentication and authorization,
	// this endpoint specifies a URI at which a client key may be signed by the actor's key for a time window to
	// act on behalf of the actor in interacting with foreign servers.
	SharedInbox as.Item `jsonld:"sharedInbox,omitempty"`
}

// Actor is generally one of the ActivityStreams actor Types, but they don't have to be.
// For example, a Profile object might be used as an actor, or a type from an ActivityStreams extension.
// Actors are retrieved like any other Object in ActivityPub.
// Like other ActivityStreams objects, actors have an id, which is a URI.
type actor struct {
	as.Parent
	// A reference to an [ActivityStreams] OrderedCollection comprised of all the messages received by the actor;
	// see 5.2 Inbox.
	Inbox as.Item `jsonld:"inbox,omitempty"`
	// An [ActivityStreams] OrderedCollection comprised of all the messages produced by the actor;
	// see 5.1 Outbox.
	Outbox as.Item `jsonld:"outbox,omitempty"`
	// A link to an [ActivityStreams] collection of the actors that this actor is following;
	// see 5.4 Following Collection
	Following as.Item `jsonld:"following,omitempty"`
	// A link to an [ActivityStreams] collection of the actors that follow this actor;
	// see 5.3 Followers Collection.
	Followers as.Item `jsonld:"followers,omitempty"`
	// A link to an [ActivityStreams] collection of the actors that follow this actor;
	// see 5.3 Followers Collection.
	Liked as.Item `jsonld:"liked,omitempty"`
	// A short username which may be used to refer to the actor, with no uniqueness guarantees.
	PreferredUsername as.NaturalLanguageValues `jsonld:"preferredUsername,omitempty,collapsible"`
	// A json object which maps additional (typically server/domain-wide) endpoints which may be useful either
	// for this actor or someone referencing this actor.
	// This mapping may be nested inside the actor document as the value or may be a link
	// to a JSON-LD document with these properties.
	Endpoints Endpoints `jsonld:"endpoints,omitempty"`
	// A list of supplementary Collections which may be of interest.
	Streams []as.ItemCollection `jsonld:"streams,omitempty"`
}

type (
	// Application describes a software application.
	Application = actor

	// Group represents a formal or informal collective of Actors.
	Group = actor

	// Organization represents an organization.
	Organization = actor

	// Person represents an individual person.
	Person = actor

	// Service represents a service of any kind.
	Service = actor
)

// actorNew initializes an actor type actor
func actorNew(id as.ObjectID, typ as.ActivityVocabularyType) *actor {
	if !as.ValidActorType(typ) {
		typ = as.ActorType
	}

	a := actor{Parent: as.Object {ID: id, Type: typ}}
	a.Name = as.NaturalLanguageValuesNew()
	a.Content = as.NaturalLanguageValuesNew()
	a.Summary = as.NaturalLanguageValuesNew()
	in := as.OrderedCollectionNew(as.ObjectID("test-inbox"))
	out := as.OrderedCollectionNew(as.ObjectID("test-outbox"))
	liked := as.OrderedCollectionNew(as.ObjectID("test-liked"))

	a.Inbox = in
	a.Outbox = out
	a.Liked = liked
	a.PreferredUsername = as.NaturalLanguageValuesNew()

	return &a
}

// ApplicationNew initializes an Application type actor
func ApplicationNew(id as.ObjectID) *Application {
	a := actorNew(id, as.ApplicationType)
	o := Application(*a)
	return &o
}

// GroupNew initializes a Group type actor
func GroupNew(id as.ObjectID) *Group {
	a := actorNew(id, as.GroupType)
	o := Group(*a)
	return &o
}

// OrganizationNew initializes an Organization type actor
func OrganizationNew(id as.ObjectID) *Organization {
	a := actorNew(id, as.OrganizationType)
	o := Organization(*a)
	return &o
}

// PersonNew initializes a Person type actor
func PersonNew(id as.ObjectID) *Person {
	a := actorNew(id, as.PersonType)
	o := Person(*a)
	return &o
}

// ServiceNew initializes a Service type actor
func ServiceNew(id as.ObjectID) *Service {
	a := actorNew(id, as.ServiceType)
	o := Service(*a)
	return &o
}

func (a *actor)UnmarshalJSON(data []byte) error {
	as.ItemTyperFunc = JSONGetItemByType

	a.Parent.UnmarshalJSON(data)
	a.PreferredUsername = as.JSONGetNaturalLanguageField(data, "preferredUsername")
	a.Followers = as.JSONGetItem(data, "followers")
	a.Following = as.JSONGetItem(data, "following")
	a.Inbox = as.JSONGetItem(data, "inbox")
	a.Outbox = as.JSONGetItem(data, "outbox")
	a.Liked = as.JSONGetItem(data, "liked")
	// TODO(marius): Endpoints need their own UnmarshalJSON
	//a.Endpoints = as.JSONGetItems(data, "endpoints")
	// TODO(marius): Streams needs custom unmarshalling
	//a.Streams = as.JSONGetItems(data, "streams")
	return nil
}