package activitypub

// Item struct
type Item = ObjectOrLink

const (
	// EmptyIRI represents a zero length IRI
	EmptyIRI IRI = ""
	// NilIRI represents by convention an IRI which is nil
	// Its use is mostly to check if a property of an ActivityPub Item is nil
	NilIRI IRI = "-"

	// EmptyID represents a zero length ID
	EmptyID = EmptyIRI
	// NilID represents by convention an ID which is nil, see details of NilIRI
	NilID = NilIRI
)

// ItemsEqual checks if it and with Items are equal
func ItemsEqual(it, with Item) bool {
	if IsNil(it) || IsNil(with) {
		return with == it
	}
	result := false
	if it.IsCollection() {
		if it.GetType() == CollectionOfItems {
			OnItemCollection(it, func(c *ItemCollection) error {
				result = c.Equals(with)
				return nil
			})
		}
		if it.GetType() == CollectionType {
			OnCollection(it, func(c *Collection) error {
				result = c.Equals(with)
				return nil
			})
		}
		if it.GetType() == OrderedCollectionType {
			OnOrderedCollection(it, func(c *OrderedCollection) error {
				result = c.Equals(with)
				return nil
			})
		}
		if it.GetType() == CollectionPageType {
			OnCollectionPage(it, func(c *CollectionPage) error {
				result = c.Equals(with)
				return nil
			})
		}
		if it.GetType() == OrderedCollectionPageType {
			OnOrderedCollectionPage(it, func(c *OrderedCollectionPage) error {
				result = c.Equals(with)
				return nil
			})
		}
	}
	if it.IsObject() {
		if ActivityTypes.Contains(with.GetType()) {
			OnActivity(it, func(i *Activity) error {
				result = i.Equals(with)
				return nil
			})
		} else if ActorTypes.Contains(with.GetType()) {
			OnActor(it, func(i *Actor) error {
				result = i.Equals(with)
				return nil
			})
		} else {
			OnObject(it, func(i *Object) error {
				result = i.Equals(with)
				return nil
			})
		}
	}
	if with.IsLink() {
		result = with.GetLink().Equals(it.GetLink(), false)
	}
	return result
}

// IsItemCollection returns if the current Item interface holds a Collection
func IsItemCollection(it Item) bool {
	_, ok := it.(ItemCollection)
	_, okP := it.(*ItemCollection)
	return ok || okP
}

// IsIRI returns if the current Item interface holds an IRI
func IsIRI(it Item) bool {
	_, okV := it.(IRI)
	_, okP := it.(*IRI)
	return okV || okP
}

// IsObject returns if the current Item interface holds an IRI
func IsObject(it Item) bool {
	switch it.(type) {
	case Actor, *Actor,
		Object, *Object, Profile, *Profile, Place, *Place, Relationship, *Relationship, Tombstone, *Tombstone,
		Activity, *Activity, IntransitiveActivity, *IntransitiveActivity, Question, *Question,
		Collection, *Collection, CollectionPage, *CollectionPage,
		OrderedCollection, *OrderedCollection, OrderedCollectionPage, *OrderedCollectionPage:
		return true
	default:
		return false
	}
}

// IsNil checks if the object matching an ObjectOrLink interface is nil
func IsNil(it Item) bool {
	if it == nil {
		return true
	}
	// This is the default if the argument can't be casted to Object, as is the case for an ItemCollection
	isNil := false
	if IsItemCollection(it) {
		OnItemCollection(it, func(c *ItemCollection) error {
			isNil = c == nil
			return nil
		})
	} else {
		OnObject(it, func(o *Object) error {
			isNil = o == nil
			return nil
		})
	}
	return isNil
}
