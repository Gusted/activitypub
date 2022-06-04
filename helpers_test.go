package activitypub

import (
	"fmt"
	"testing"
)

func assertObjectWithTesting(fn canErrorFunc, expected Item) WithObjectFn {
	return func(p *Object) error {
		if !assertDeepEquals(fn, p, expected) {
			return fmt.Errorf("not equal")
		}
		return nil
	}
}

func TestOnObject(t *testing.T) {
	testObject := Object{
		ID: "https://example.com",
	}
	type args struct {
		it Item
		fn func(canErrorFunc, Item) WithObjectFn
	}
	tests := []struct {
		name     string
		args     args
		expected Item
		wantErr  bool
	}{
		{
			name:     "single",
			args:     args{testObject, assertObjectWithTesting},
			expected: &testObject,
			wantErr:  false,
		},
		{
			name:     "single fails",
			args:     args{Object{ID: "https://not-equals"}, assertObjectWithTesting},
			expected: &testObject,
			wantErr:  true,
		},
		{
			name:     "collectionOfObjects",
			args:     args{ItemCollection{testObject, testObject}, assertObjectWithTesting},
			expected: &testObject,
			wantErr:  false,
		},
		{
			name:     "collectionOfObjects fails",
			args:     args{ItemCollection{testObject, Object{ID: "https://not-equals"}}, assertObjectWithTesting},
			expected: &testObject,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		var logFn canErrorFunc
		if tt.wantErr {
			logFn = t.Logf
		} else {
			logFn = t.Errorf
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := OnObject(tt.args.it, tt.args.fn(logFn, tt.expected)); (err != nil) != tt.wantErr {
				t.Errorf("OnObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func assertActivityWithTesting(fn canErrorFunc, expected Item) WithActivityFn {
	return func(p *Activity) error {
		if !assertDeepEquals(fn, p, expected) {
			return fmt.Errorf("not equal")
		}
		return nil
	}
}

func TestOnActivity(t *testing.T) {
	testActivity := Activity{
		ID: "https://example.com",
	}
	type args struct {
		it Item
		fn func(canErrorFunc, Item) WithActivityFn
	}
	tests := []struct {
		name     string
		args     args
		expected Item
		wantErr  bool
	}{
		{
			name:     "single",
			args:     args{testActivity, assertActivityWithTesting},
			expected: &testActivity,
			wantErr:  false,
		},
		{
			name:     "single fails",
			args:     args{Activity{ID: "https://not-equals"}, assertActivityWithTesting},
			expected: &testActivity,
			wantErr:  true,
		},
		{
			name:     "collectionOfActivitys",
			args:     args{ItemCollection{testActivity, testActivity}, assertActivityWithTesting},
			expected: &testActivity,
			wantErr:  false,
		},
		{
			name:     "collectionOfActivitys fails",
			args:     args{ItemCollection{testActivity, Activity{ID: "https://not-equals"}}, assertActivityWithTesting},
			expected: &testActivity,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		var logFn canErrorFunc
		if tt.wantErr {
			logFn = t.Logf
		} else {
			logFn = t.Errorf
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := OnActivity(tt.args.it, tt.args.fn(logFn, tt.expected)); (err != nil) != tt.wantErr {
				t.Errorf("OnActivity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func assertIntransitiveActivityWithTesting(fn canErrorFunc, expected Item) WithIntransitiveActivityFn {
	return func(p *IntransitiveActivity) error {
		if !assertDeepEquals(fn, p, expected) {
			return fmt.Errorf("not equal")
		}
		return nil
	}
}

func TestOnIntransitiveActivity(t *testing.T) {
	testIntransitiveActivity := IntransitiveActivity{
		ID: "https://example.com",
	}
	type args struct {
		it Item
		fn func(canErrorFunc, Item) WithIntransitiveActivityFn
	}
	tests := []struct {
		name     string
		args     args
		expected Item
		wantErr  bool
	}{
		{
			name:     "single",
			args:     args{testIntransitiveActivity, assertIntransitiveActivityWithTesting},
			expected: &testIntransitiveActivity,
			wantErr:  false,
		},
		{
			name:     "single fails",
			args:     args{IntransitiveActivity{ID: "https://not-equals"}, assertIntransitiveActivityWithTesting},
			expected: &testIntransitiveActivity,
			wantErr:  true,
		},
		{
			name:     "collectionOfIntransitiveActivitys",
			args:     args{ItemCollection{testIntransitiveActivity, testIntransitiveActivity}, assertIntransitiveActivityWithTesting},
			expected: &testIntransitiveActivity,
			wantErr:  false,
		},
		{
			name:     "collectionOfIntransitiveActivitys fails",
			args:     args{ItemCollection{testIntransitiveActivity, IntransitiveActivity{ID: "https://not-equals"}}, assertIntransitiveActivityWithTesting},
			expected: &testIntransitiveActivity,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		var logFn canErrorFunc
		if tt.wantErr {
			logFn = t.Logf
		} else {
			logFn = t.Errorf
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := OnIntransitiveActivity(tt.args.it, tt.args.fn(logFn, tt.expected)); (err != nil) != tt.wantErr {
				t.Errorf("OnIntransitiveActivity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func assertQuestionWithTesting(fn canErrorFunc, expected Item) WithQuestionFn {
	return func(p *Question) error {
		if !assertDeepEquals(fn, p, expected) {
			return fmt.Errorf("not equal")
		}
		return nil
	}
}

func TestOnQuestion(t *testing.T) {
	testQuestion := Question{
		ID: "https://example.com",
	}
	type args struct {
		it Item
		fn func(canErrorFunc, Item) WithQuestionFn
	}
	tests := []struct {
		name     string
		args     args
		expected Item
		wantErr  bool
	}{
		{
			name:     "single",
			args:     args{testQuestion, assertQuestionWithTesting},
			expected: &testQuestion,
			wantErr:  false,
		},
		{
			name:     "single fails",
			args:     args{Question{ID: "https://not-equals"}, assertQuestionWithTesting},
			expected: &testQuestion,
			wantErr:  true,
		},
		{
			name:     "collectionOfQuestions",
			args:     args{ItemCollection{testQuestion, testQuestion}, assertQuestionWithTesting},
			expected: &testQuestion,
			wantErr:  false,
		},
		{
			name:     "collectionOfQuestions fails",
			args:     args{ItemCollection{testQuestion, Question{ID: "https://not-equals"}}, assertQuestionWithTesting},
			expected: &testQuestion,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		var logFn canErrorFunc
		if tt.wantErr {
			logFn = t.Logf
		} else {
			logFn = t.Errorf
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := OnQuestion(tt.args.it, tt.args.fn(logFn, tt.expected)); (err != nil) != tt.wantErr {
				t.Errorf("OnQuestion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOnCollection(t *testing.T) {
	t.Skipf("TODO")
}

func TestOnCollectionPage(t *testing.T) {
	t.Skipf("TODO")
}

func TestOnOrderedCollectionPage(t *testing.T) {
	t.Skipf("TODO")
}

type args[T Objects] struct {
	it T
	fn func(fn canErrorFunc, expected T) func(*T) error
}

type testPair[T Objects] struct {
	name     string
	args     args[T]
	expected T
	wantErr  bool
}

func assert[T Objects](fn canErrorFunc, expected T) func(*T) error {
	return func(p *T) error {
		if !assertDeepEquals(fn, *p, expected) {
			return fmt.Errorf("not equal")
		}
		return nil
	}
}

func TestOn(t *testing.T) {
	var tests = []testPair[Object]{
		{
			name:     "single object",
			args:     args[Object]{Object{ID: "https://example.com"}, assert[Object]},
			expected: Object{ID: "https://example.com"},
			wantErr:  false,
		},
		{
			name:     "single image",
			args:     args[Image]{Image{ID: "http://example.com"}, assert[Image]},
			expected: Image{ID: "http://example.com"},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		var logFn canErrorFunc
		if tt.wantErr {
			logFn = t.Logf
		} else {
			logFn = t.Errorf
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := On(tt.args.it, tt.args.fn(logFn, tt.expected)); (err != nil) != tt.wantErr {
				t.Errorf("On[%T]() error = %v, wantErr %v", tt.args.it, err, tt.wantErr)
			}
		})
	}
}

var fnObj = func(_ *Object) error {
	return nil
}
var fnAct = func(_ *Actor) error {
	return nil
}
var fnA = func(_ *Activity) error {
	return nil
}
var fnIA = func(_ *IntransitiveActivity) error {
	return nil
}

func Benchmark_OnT_vs_On_T(b *testing.B) {
	var it Item
	b.Run("OnObject", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			OnObject(it, fnObj)
		}
	})
	b.Run("On_T_Object", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			On[Object](it, fnObj)
		}
	})
	b.Run("OnActor", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			OnActor(it, fnAct)
		}
	})
	b.Run("On_T_Actor", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			On[Actor](it, fnAct)
		}
	})
	b.Run("OnActivity", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			OnActivity(it, fnA)
		}
	})
	b.Run("On_T_Activity", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			On[Activity](it, fnA)
		}
	})
	b.Run("OnIntransitiveActivity", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			OnIntransitiveActivity(it, fnIA)
		}
	})
	b.Run("On_T_IntransitiveActivity", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			On[IntransitiveActivity](it, fnIA)
		}
	})
}
