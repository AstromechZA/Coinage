// Package assert provides a handful of useful assertion primitives for basic Golang unit tests.
// This is not meant to replace packages like testify but is meant to be used in small library code that don't want
// to rely on a large 3rd party dependency just for a handful of assert tests.
//
// Copy and past this file into your code base as internal/assert/assert.go and use it in your testing.
package assert

import (
	"reflect"
	"testing"
)

func Equal(t *testing.T, a, b interface{}) {
	t.Helper()
	Equalf(t, a, b, "%v (%T) != %v (%T)", a, a, b, b)
}

func Equalf(t *testing.T, a interface{}, b interface{}, msg string, args ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf(msg, args...)
	}
}

func NotEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	NotEqualf(t, a, b, "%v (%T) == %v (%T)", a, a, b, b)
}

func NotEqualf(t *testing.T, a interface{}, b interface{}, msg string, args ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		t.Fatalf(msg, args...)
	}
}

func ShouldEqual(t *testing.T, a, b interface{}) bool {
	t.Helper()
	return ShouldEqualf(t, a, b, "%v (%T) != %v (%T)", a, a, b, b)
}

func ShouldEqualf(t *testing.T, a interface{}, b interface{}, msg string, args ...interface{}) bool {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf(msg, args...)
		return false
	}
	return true
}

func ShouldNotEqual(t *testing.T, a, b interface{}) bool {
	t.Helper()
	return ShouldNotEqualf(t, a, b, "%v (%T) == %v (%T)", a, a, b, b)
}

func ShouldNotEqualf(t *testing.T, a interface{}, b interface{}, msg string, args ...interface{}) bool {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Errorf(msg, args...)
		return false
	}
	return true
}

func True(t *testing.T, a interface{}) {
	t.Helper()
	Equal(t, a, true)
}

func ShouldBeTrue(t *testing.T, a interface{}) bool {
	t.Helper()
	return ShouldEqual(t, a, true)
}

func False(t *testing.T, a interface{}) {
	t.Helper()
	Equal(t, a, false)
}

func ShouldBeFalse(t *testing.T, a interface{}) bool {
	t.Helper()
	return ShouldEqual(t, a, false)
}
