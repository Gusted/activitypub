package handler

import (
	"fmt"
	as "github.com/go-ap/activitystreams"
	j "github.com/go-ap/jsonld"
	"net/http"
)

// ActivityHandlerFn is the type that we're using to represent handlers that process requests containing
// an ActivityStreams Activity. It needs to implement the http.Handler interface.
//
// It is considered that following the execution of the handler, we return a pair formed of a HTTP status together with
//  an IRI representing a new Object - in the case of transitive activities that had a side effect, or
//  an error.
// In the case of intransitive activities the iri will always be empty.
type ActivityHandlerFn func(http.ResponseWriter, *http.Request) (as.IRI, int, error)

// ValidMethod validates if the current handler can process the current request
func (a ActivityHandlerFn) ValidMethod(r *http.Request) bool {
	return r.Method != http.MethodPost
}

// ValidateRequest validates if the current handler can process the current request
func (a ActivityHandlerFn) ValidateRequest(r *http.Request) (int, error) {
	if !a.ValidMethod(r) {
		return http.StatusNotAcceptable, fmt.Errorf("invalid HTTP method %s", r.Method)
	}
	return http.StatusOK, nil
}

// ServeHTTP implements the http.Handler interface for the ActivityHandlerFn type
func (a ActivityHandlerFn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dat []byte
	var iri as.IRI
	var err error
	var status = http.StatusOK

	if status, err = a.ValidateRequest(r); err != nil {
		dat = []byte(err.Error())
	}

	if iri, status, err = a(w, r); err != nil {
		status = http.StatusInternalServerError
		dat = []byte(err.Error())
	} else {
		dat = []byte("OK")
	}

	w.WriteHeader(status)
	if len(iri) > 0 {
		w.Header().Set("Location", iri.String())
	}
	w.Write(dat)
}

type ClientHandler interface {
	ValidateRequest(r *http.Request) (int, error)
}

// CollectionHandlerFn is the type that we're using to represent handlers that will return ActivityStreams
// Collection or OrderedCollection objects. It needs to implement the http.Handler interface.
type CollectionHandlerFn func(http.ResponseWriter, *http.Request) (as.CollectionInterface, error)

// ValidMethod validates if the current handler can process the current request
func (c CollectionHandlerFn) ValidMethod(r *http.Request) bool {
	return r.Method == http.MethodGet || r.Method == http.MethodHead
}

// ValidateRequest validates if the current handler can process the current request
func (c CollectionHandlerFn) ValidateRequest(r *http.Request) (int, error) {
	if !c.ValidMethod(r) {
		return http.StatusMethodNotAllowed, fmt.Errorf("invalid HTTP method %s", r.Method)
	}
	return http.StatusOK, nil
}

// ServeHTTP implements the http.Handler interface for the CollectionHandlerFn type
func (c CollectionHandlerFn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dat []byte

	var status = http.StatusOK
	var err error

	if status, err = c.ValidateRequest(r); err != nil {
		dat = []byte(err.Error())
	} else {
		if col, err := c(w, r); err != nil {
			dat = []byte(err.Error())
		} else {
			if dat, err = j.WithContext(j.IRI(as.ActivityBaseURI)).Marshal(col); err != nil {
				status = http.StatusInternalServerError
				dat = []byte(err.Error())
			}
		}
	}

	w.WriteHeader(status)
	if r.Method == http.MethodGet {
		w.Write(dat)
	}
}

// ItemHandlerFn is the type that we're using to represent handlers that return ActivityStreams
// objects. It needs to implement the http.Handler interface
type ItemHandlerFn func(http.ResponseWriter, *http.Request) (as.Item, error)

// ValidMethod validates if the current handler can process the current request
func (i ItemHandlerFn) ValidMethod(r *http.Request) bool {
	return r.Method != http.MethodGet && r.Method != http.MethodHead
}

// ValidateRequest validates if the current handler can process the current request
func (i ItemHandlerFn) ValidateRequest(r *http.Request) (int, error) {
	if !i.ValidMethod(r) {
		return http.StatusMethodNotAllowed, fmt.Errorf("invalid HTTP method %s", r.Method)
	}
	return http.StatusOK, nil
}

// ServeHTTP implements the http.Handler interface for the ItemHandlerFn type
func (i ItemHandlerFn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var dat []byte
	var status = http.StatusOK
	var err error

	if status, err = i.ValidateRequest(r); err != nil {
		dat = []byte(err.Error())
	} else {
		if it, err := i(w, r); err != nil {
			status = http.StatusInternalServerError
			dat = []byte(err.Error())
		} else {
			if dat, err = j.WithContext(j.IRI(as.ActivityBaseURI)).Marshal(it); err != nil {
				status = http.StatusInternalServerError
				dat = []byte(err.Error())
			}
		}
	}

	w.WriteHeader(status)
	if r.Method == http.MethodGet {
		w.Write(dat)
	}
}
