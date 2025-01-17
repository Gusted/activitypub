package activitypub

import (
	"encoding"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/valyala/fastjson"
)

var (
	apUnmarshalerType   = reflect.TypeOf(new(Item)).Elem()
	unmarshalerType     = reflect.TypeOf(new(json.Unmarshaler)).Elem()
	textUnmarshalerType = reflect.TypeOf(new(encoding.TextUnmarshaler)).Elem()
)

// ItemTyperFunc will return an instance of a struct that implements activitystreams.Item
// The default for this package is GetItemByType but can be overwritten
var ItemTyperFunc TyperFn = GetItemByType

// TyperFn is the type of the function which returns an activitystreams.Item struct instance
// for a specific ActivityVocabularyType
type TyperFn func(ActivityVocabularyType) (Item, error)

func JSONGetID(val *fastjson.Value) ID {
	i := val.Get("id").GetStringBytes()
	return ID(i)
}

func JSONGetType(val *fastjson.Value) ActivityVocabularyType {
	t := val.Get("type").GetStringBytes()
	return ActivityVocabularyType(t)
}

func JSONGetMimeType(val *fastjson.Value, prop string) MimeType {
	if !val.Exists(prop) {
		return ""
	}
	t := val.GetStringBytes(prop)
	return MimeType(t)
}

func JSONGetInt(val *fastjson.Value, prop string) int64 {
	if !val.Exists(prop) {
		return 0
	}
	i := val.Get(prop).GetInt64()
	return i
}

func JSONGetFloat(val *fastjson.Value, prop string) float64 {
	if !val.Exists(prop) {
		return 0.0
	}
	f := val.Get(prop).GetFloat64()
	return f
}

func JSONGetString(val *fastjson.Value, prop string) string {
	if !val.Exists(prop) {
		return ""
	}
	s := val.Get(prop).GetStringBytes()
	return string(s)
}

func JSONGetBytes(val *fastjson.Value, prop string) []byte {
	if !val.Exists(prop) {
		return nil
	}
	s := val.Get(prop).GetStringBytes()
	return s
}

func JSONGetBoolean(val *fastjson.Value, prop string) bool {
	if !val.Exists(prop) {
		return false
	}
	t, _ := val.Get(prop).Bool()
	return t
}

func JSONGetNaturalLanguageField(val *fastjson.Value, prop string) NaturalLanguageValues {
	n := NaturalLanguageValues{}
	if val == nil {
		return n
	}
	v := val.Get(prop)
	if v == nil {
		return nil
	}
	switch v.Type() {
	case fastjson.TypeObject:
		ob, _ := v.Object()
		ob.Visit(func(key []byte, v *fastjson.Value) {
			l := LangRefValue{}
			l.Ref = LangRef(key)
			if err := l.Value.UnmarshalJSON(v.GetStringBytes()); err == nil {
				if l.Ref != NilLangRef || len(l.Value) > 0 {
					n = append(n, l)
				}
			}
		})
	case fastjson.TypeString:
		l := LangRefValue{}
		if err := l.UnmarshalJSON(v.GetStringBytes()); err == nil {
			n = append(n, l)
		}
	}

	return n
}

func JSONGetTime(val *fastjson.Value, prop string) time.Time {
	t := time.Time{}
	if val == nil {
		return t
	}

	if str := val.Get(prop).GetStringBytes(); len(str) > 0 {
		t.UnmarshalText(str)
		return t.UTC()
	}
	return t
}

func JSONGetDuration(val *fastjson.Value, prop string) time.Duration {
	if str := val.Get(prop).GetStringBytes(); len(str) > 0 {
		// TODO(marius): this needs to be replaced to be compatible with xsd:duration
		d, _ := time.ParseDuration(string(str))
		return d
	}
	return 0
}

func JSONGetPublicKey(val *fastjson.Value, prop string) PublicKey {
	key := PublicKey{}
	if val == nil {
		return key
	}
	val = val.Get(prop)
	if val == nil {
		return key
	}
	loadPublicKey(val, &key)
	return key
}

func itemsFn(val *fastjson.Value) (Item, error) {
	if val.Type() == fastjson.TypeArray {
		it := val.GetArray()
		if len(it) == 1 {
			return itemFn(it[0])
		}
		items := make(ItemCollection, 0)
		for _, v := range it {
			if it, _ := itemFn(v); it != nil {
				items.Append(it)
			}
		}
		return items, nil
	}
	return itemFn(val)
}

func itemFn(val *fastjson.Value) (Item, error) {
	typ := JSONGetType(val)
	if typ == "" && val.Type() == fastjson.TypeString {
		// try to see if it's an IRI
		if i, ok := asIRI(val); ok {
			return i, nil
		}
	}
	i, err := ItemTyperFunc(typ)
	if err != nil || IsNil(i) {
		return nil, nil
	}

	switch typ {
	case "":
		// NOTE(marius): this handles Tags, that don't have types
		fallthrough
	case ObjectType, ArticleType, AudioType, DocumentType, EventType, ImageType, NoteType, PageType, VideoType:
		err = OnObject(i, func(ob *Object) error {
			return loadObject(val, ob)
		})
	case LinkType, MentionType:
		err = OnLink(i, func(l *Link) error {
			return loadLink(val, l)
		})
	case ActivityType, AcceptType, AddType, AnnounceType, BlockType, CreateType, DeleteType, DislikeType,
		FlagType, FollowType, IgnoreType, InviteType, JoinType, LeaveType, LikeType, ListenType, MoveType, OfferType,
		RejectType, ReadType, RemoveType, TentativeRejectType, TentativeAcceptType, UndoType, UpdateType, ViewType:
		err = OnActivity(i, func(act *Activity) error {
			return loadActivity(val, act)
		})
	case IntransitiveActivityType, ArriveType, TravelType:
		err = OnIntransitiveActivity(i, func(act *IntransitiveActivity) error {
			return loadIntransitiveActivity(val, act)
		})
	case ActorType, ApplicationType, GroupType, OrganizationType, PersonType, ServiceType:
		err = OnActor(i, func(a *Actor) error {
			return loadActor(val, a)
		})
	case CollectionType:
		err = OnCollection(i, func(c *Collection) error {
			return loadCollection(val, c)
		})
	case OrderedCollectionType:
		err = OnOrderedCollection(i, func(c *OrderedCollection) error {
			return loadOrderedCollection(val, c)
		})
	case CollectionPageType:
		err = OnCollectionPage(i, func(p *CollectionPage) error {
			return loadCollectionPage(val, p)
		})
	case OrderedCollectionPageType:
		err = OnOrderedCollectionPage(i, func(p *OrderedCollectionPage) error {
			return loadOrderedCollectionPage(val, p)
		})
	case PlaceType:
		err = OnPlace(i, func(p *Place) error {
			return loadPlace(val, p)
		})
	case ProfileType:
		err = OnProfile(i, func(p *Profile) error {
			return loadProfile(val, p)
		})
	case RelationshipType:
		err = OnRelationship(i, func(r *Relationship) error {
			return loadRelationship(val, r)
		})
	case TombstoneType:
		err = OnTombstone(i, func(t *Tombstone) error {
			return loadTombstone(val, t)
		})
	case QuestionType:
		err = OnQuestion(i, func(q *Question) error {
			return loadQuestion(val, q)
		})
	}

	if !NotEmpty(i) {
		return nil, nil
	}
	return i, err
}

func JSONUnmarshalToItem(val *fastjson.Value) Item {
	if ItemTyperFunc == nil {
		return nil
	}
	var (
		i   Item
		err error
	)
	switch val.Type() {
	case fastjson.TypeArray:
		i, err = itemsFn(val)
	case fastjson.TypeObject:
		i, err = itemFn(val)
	case fastjson.TypeString:
		if iri, ok := asIRI(val); ok {
			// try to see if it's an IRI
			i = iri
		}
	}
	if err != nil {
		return nil
	}
	return i
}

func asIRI(val *fastjson.Value) (IRI, bool) {
	if val == nil {
		return NilIRI, true
	}
	s := strings.Trim(val.String(), `"`)
	u, err := url.ParseRequestURI(s)
	if err == nil && len(u.Scheme) > 0 && len(u.Host) > 0 {
		// try to see if it's an IRI
		return IRI(s), true
	}
	return EmptyIRI, false
}

func JSONGetItem(val *fastjson.Value, prop string) Item {
	if val == nil {
		return nil
	}
	if val = val.Get(prop); val == nil {
		return nil
	}
	switch val.Type() {
	case fastjson.TypeString:
		if i, ok := asIRI(val); ok {
			// try to see if it's an IRI
			return i
		}
	case fastjson.TypeArray:
		it, _ := itemsFn(val)
		return it
	case fastjson.TypeObject:
		it, _ := itemFn(val)
		return it
	case fastjson.TypeNumber:
		fallthrough
	case fastjson.TypeNull:
		fallthrough
	default:
		return nil
	}
	return nil
}

func JSONGetURIItem(val *fastjson.Value, prop string) Item {
	if val == nil {
		return nil
	}
	if val = val.Get(prop); val == nil {
		return nil
	}
	switch val.Type() {
	case fastjson.TypeObject:
		if it, _ := itemFn(val); it != nil {
			return it
		}
	case fastjson.TypeArray:
		if it, _ := itemsFn(val); it != nil {
			return it
		}
	case fastjson.TypeString:
		return IRI(val.GetStringBytes())
	}

	return nil
}

func JSONGetItems(val *fastjson.Value, prop string) ItemCollection {
	if val == nil {
		return nil
	}
	if val = val.Get(prop); val == nil {
		return nil
	}

	it := make(ItemCollection, 0)
	switch val.Type() {
	case fastjson.TypeArray:
		for _, v := range val.GetArray() {
			if i, _ := itemFn(v); i != nil {
				it.Append(i)
			}
		}
	case fastjson.TypeObject:
		if i := JSONGetItem(val, prop); i != nil {
			it.Append(i)
		}
	case fastjson.TypeString:
		if iri := val.GetStringBytes(); len(iri) > 0 {
			it.Append(IRI(iri))
		}
	}
	if len(it) == 0 {
		return nil
	}
	return it
}

func JSONGetLangRefField(val *fastjson.Value, prop string) LangRef {
	s := val.Get(prop).GetStringBytes()
	return LangRef(s)
}

func JSONGetIRI(val *fastjson.Value, prop string) IRI {
	s := val.Get(prop).GetStringBytes()
	return IRI(s)
}

// UnmarshalJSON tries to detect the type of the object in the json data and then outputs a matching
// ActivityStreams object, if possible
func UnmarshalJSON(data []byte) (Item, error) {
	if ItemTyperFunc == nil {
		ItemTyperFunc = GetItemByType
	}
	p := fastjson.Parser{}
	val, err := p.ParseBytes(data)
	if err != nil {
		return nil, err
	}
	return JSONUnmarshalToItem(val), nil
}

func GetItemByType(typ ActivityVocabularyType) (Item, error) {
	switch typ {
	case ObjectType, ArticleType, AudioType, DocumentType, EventType, ImageType, NoteType, PageType, VideoType:
		return ObjectNew(typ), nil
	case LinkType, MentionType:
		return &Link{Type: typ}, nil
	case ActivityType, AcceptType, AddType, AnnounceType, BlockType, CreateType, DeleteType, DislikeType,
		FlagType, FollowType, IgnoreType, InviteType, JoinType, LeaveType, LikeType, ListenType, MoveType, OfferType,
		RejectType, ReadType, RemoveType, TentativeRejectType, TentativeAcceptType, UndoType, UpdateType, ViewType:
		return &Activity{Type: typ}, nil
	case IntransitiveActivityType, ArriveType, TravelType:
		return &IntransitiveActivity{Type: typ}, nil
	case ActorType, ApplicationType, GroupType, OrganizationType, PersonType, ServiceType:
		return &Actor{Type: typ}, nil
	case CollectionType:
		return &Collection{Type: typ}, nil
	case OrderedCollectionType:
		return &OrderedCollection{Type: typ}, nil
	case CollectionPageType:
		return &CollectionPage{Type: typ}, nil
	case OrderedCollectionPageType:
		return &OrderedCollectionPage{Type: typ}, nil
	case PlaceType:
		return &Place{Type: typ}, nil
	case ProfileType:
		return &Profile{Type: typ}, nil
	case RelationshipType:
		return &Relationship{Type: typ}, nil
	case TombstoneType:
		return &Tombstone{Type: typ}, nil
	case QuestionType:
		return &Question{Type: typ}, nil
	case "":
		// when no type is available use a plain Object
		return &Object{}, nil
	}
	return nil, fmt.Errorf("empty ActivityStreams type")
}

func JSONGetActorEndpoints(val *fastjson.Value, prop string) *Endpoints {
	str := val.Get(prop).GetStringBytes()

	var e *Endpoints
	if len(str) > 0 {
		e = &Endpoints{}
		e.UnmarshalJSON(str)
	}

	return e
}

func loadObject(val *fastjson.Value, o *Object) error {
	if ItemTyperFunc == nil {
		ItemTyperFunc = GetItemByType
	}
	o.ID = JSONGetID(val)
	o.Type = JSONGetType(val)
	o.Name = JSONGetNaturalLanguageField(val, "name")
	o.Content = JSONGetNaturalLanguageField(val, "content")
	o.Summary = JSONGetNaturalLanguageField(val, "summary")
	o.Context = JSONGetItem(val, "context")
	o.URL = JSONGetURIItem(val, "url")
	o.MediaType = JSONGetMimeType(val, "mediaType")
	o.Generator = JSONGetItem(val, "generator")
	o.AttributedTo = JSONGetItem(val, "attributedTo")
	o.Attachment = JSONGetItem(val, "attachment")
	o.Location = JSONGetItem(val, "location")
	o.Published = JSONGetTime(val, "published")
	o.StartTime = JSONGetTime(val, "startTime")
	o.EndTime = JSONGetTime(val, "endTime")
	o.Duration = JSONGetDuration(val, "duration")
	o.Icon = JSONGetItem(val, "icon")
	o.Preview = JSONGetItem(val, "preview")
	o.Image = JSONGetItem(val, "image")
	o.Updated = JSONGetTime(val, "updated")
	o.InReplyTo = JSONGetItem(val, "inReplyTo")
	o.To = JSONGetItems(val, "to")
	o.Audience = JSONGetItems(val, "audience")
	o.Bto = JSONGetItems(val, "bto")
	o.CC = JSONGetItems(val, "cc")
	o.BCC = JSONGetItems(val, "bcc")
	o.Replies = JSONGetItem(val, "replies")
	o.Tag = JSONGetItems(val, "tag")
	o.Likes = JSONGetItem(val, "likes")
	o.Shares = JSONGetItem(val, "shares")
	o.Source = GetAPSource(val)
	return nil
}

func loadIntransitiveActivity(val *fastjson.Value, i *IntransitiveActivity) error {
	OnObject(i, func(o *Object) error {
		return loadObject(val, o)
	})
	i.Actor = JSONGetItem(val, "actor")
	i.Target = JSONGetItem(val, "target")
	i.Result = JSONGetItem(val, "result")
	i.Origin = JSONGetItem(val, "origin")
	i.Instrument = JSONGetItem(val, "instrument")
	return nil
}

func loadActivity(val *fastjson.Value, a *Activity) error {
	OnIntransitiveActivity(a, func(i *IntransitiveActivity) error {
		return loadIntransitiveActivity(val, i)
	})
	a.Object = JSONGetItem(val, "object")
	return nil
}

func loadQuestion(val *fastjson.Value, q *Question) error {
	OnIntransitiveActivity(q, func(i *IntransitiveActivity) error {
		return loadIntransitiveActivity(val, i)
	})
	q.OneOf = JSONGetItem(val, "oneOf")
	q.AnyOf = JSONGetItem(val, "anyOf")
	q.Closed = JSONGetBoolean(val, "closed")
	return nil
}

func loadActor(val *fastjson.Value, a *Actor) error {
	OnObject(a, func(o *Object) error {
		return loadObject(val, o)
	})
	a.PreferredUsername = JSONGetNaturalLanguageField(val, "preferredUsername")
	a.Followers = JSONGetItem(val, "followers")
	a.Following = JSONGetItem(val, "following")
	a.Inbox = JSONGetItem(val, "inbox")
	a.Outbox = JSONGetItem(val, "outbox")
	a.Liked = JSONGetItem(val, "liked")
	a.Endpoints = JSONGetActorEndpoints(val, "endpoints")
	a.Streams = JSONGetItems(val, "streams")
	a.PublicKey = JSONGetPublicKey(val, "publicKey")
	return nil
}

func loadCollection(val *fastjson.Value, c *Collection) error {
	OnObject(c, func(o *Object) error {
		return loadObject(val, o)
	})
	c.Current = JSONGetItem(val, "current")
	c.First = JSONGetItem(val, "first")
	c.Last = JSONGetItem(val, "last")
	c.TotalItems = uint(JSONGetInt(val, "totalItems"))
	c.Items = JSONGetItems(val, "items")
	return nil
}

func loadCollectionPage(val *fastjson.Value, c *CollectionPage) error {
	OnCollection(c, func(c *Collection) error {
		return loadCollection(val, c)
	})
	c.Next = JSONGetItem(val, "next")
	c.Prev = JSONGetItem(val, "prev")
	c.PartOf = JSONGetItem(val, "partOf")
	return nil
}

func loadOrderedCollection(val *fastjson.Value, c *OrderedCollection) error {
	OnObject(c, func(o *Object) error {
		return loadObject(val, o)
	})
	c.Current = JSONGetItem(val, "current")
	c.First = JSONGetItem(val, "first")
	c.Last = JSONGetItem(val, "last")
	c.TotalItems = uint(JSONGetInt(val, "totalItems"))
	c.OrderedItems = JSONGetItems(val, "orderedItems")
	return nil
}

func loadOrderedCollectionPage(val *fastjson.Value, c *OrderedCollectionPage) error {
	OnOrderedCollection(c, func(c *OrderedCollection) error {
		return loadOrderedCollection(val, c)
	})
	c.Next = JSONGetItem(val, "next")
	c.Prev = JSONGetItem(val, "prev")
	c.PartOf = JSONGetItem(val, "partOf")
	c.StartIndex = uint(JSONGetInt(val, "startIndex"))
	return nil
}

func loadPlace(val *fastjson.Value, p *Place) error {
	OnObject(p, func(o *Object) error {
		return loadObject(val, o)
	})
	p.Accuracy = JSONGetFloat(val, "accuracy")
	p.Altitude = JSONGetFloat(val, "altitude")
	p.Latitude = JSONGetFloat(val, "latitude")
	p.Longitude = JSONGetFloat(val, "longitude")
	p.Radius = JSONGetInt(val, "radius")
	p.Units = JSONGetString(val, "units")
	return nil
}

func loadProfile(val *fastjson.Value, p *Profile) error {
	OnObject(p, func(o *Object) error {
		return loadObject(val, o)
	})
	p.Describes = JSONGetItem(val, "describes")
	return nil
}

func loadRelationship(val *fastjson.Value, r *Relationship) error {
	OnObject(r, func(o *Object) error {
		return loadObject(val, o)
	})
	r.Subject = JSONGetItem(val, "subject")
	r.Object = JSONGetItem(val, "object")
	r.Relationship = JSONGetItem(val, "relationship")
	return nil
}

func loadTombstone(val *fastjson.Value, t *Tombstone) error {
	OnObject(t, func(o *Object) error {
		return loadObject(val, o)
	})
	t.FormerType = ActivityVocabularyType(JSONGetString(val, "formerType"))
	t.Deleted = JSONGetTime(val, "deleted")
	return nil
}

func loadLink(val *fastjson.Value, l *Link) error {
	if ItemTyperFunc == nil {
		ItemTyperFunc = GetItemByType
	}
	l.ID = JSONGetID(val)
	l.Type = JSONGetType(val)
	l.MediaType = JSONGetMimeType(val, "mediaType")
	l.Preview = JSONGetItem(val, "preview")
	h := JSONGetInt(val, "height")
	if h != 0 {
		l.Height = uint(h)
	}
	w := JSONGetInt(val, "width")
	if w != 0 {
		l.Width = uint(w)
	}
	l.Name = JSONGetNaturalLanguageField(val, "name")
	hrefLang := JSONGetLangRefField(val, "hrefLang")
	if len(hrefLang) > 0 {
		l.HrefLang = hrefLang
	}
	href := JSONGetURIItem(val, "href")
	if href != nil {
		ll := href.GetLink()
		if len(ll) > 0 {
			l.Href = ll
		}
	}
	rel := JSONGetURIItem(val, "rel")
	if rel != nil {
		rr := rel.GetLink()
		if len(rr) > 0 {
			l.Rel = rr
		}
	}
	return nil
}

func loadPublicKey(val *fastjson.Value, p *PublicKey) error {
	if id := val.GetStringBytes("id"); len(id) > 0 {
		p.ID = ID(id)
	}
	if o := val.GetStringBytes("owner"); len(o) > 0 {
		p.Owner = IRI(o)
	}
	if pub := val.GetStringBytes("publicKeyPem"); len(pub) > 0 {
		p.PublicKeyPem = string(pub)
	}
	return nil
}
