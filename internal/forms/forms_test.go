package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	// r := httptest.NewRequest("POST", "/whatever", nil)
	// form := New(r.PostForm)
	// there's actually no point for creating a new request here
	// all we new is a url.Values{} to pass to New()
	form := New(url.Values{})

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	// r := httptest.NewRequest("POST", "/whatever", nil)
	// form := New(r.PostForm)
	postedData := url.Values{}
	form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	// testing errors.go Get() method
	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("shouldn't have an error, but did get one")
	}
}

func TestForm_Has(t *testing.T) {
	// r := httptest.NewRequest("POST", "/whatever", nil)
	// form := New(r.PostForm)
	postedData := url.Values{}

	form := New(postedData)
	if form.Has("whatever") {
		t.Error("form shows valid when required fields missing")
	}
	isError := form.Errors.Get("a")
	if isError != "" {
		t.Error("form shows error for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)
	if !form.Has("a") {
		t.Error("shows does not have required fields when it does")
	}

}

func TestForm_MinLength(t *testing.T) {
	// r := httptest.NewRequest("POST", "/whatever", nil)
	// form := New(r.PostForm)
	postedData := url.Values{}

	form := New(postedData)
	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}
	postedData = url.Values{}
	postedData.Add("some_field", "some_value")

	form = New(postedData)
	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows min length of 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("some_field", "some_value")

	form = New(postedData)
	form.MinLength("some_field", 10)
	if !form.Valid() {
		t.Error("does not show min length of 10 met when it is")
	}
}

func TestForm_IsEmail(t *testing.T) {
	// r := httptest.NewRequest("POST", "/whatever", nil)
	// form := New(r.PostForm)
	postedData := url.Values{}

	form := New(postedData)
	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "cadcsd@gmail.com")

	form = New(postedData)
	form.IsEmail("email")

	if !form.Valid() {
		t.Error("Got invalid email when should have been valid")
	}

	postedData = url.Values{}
	postedData.Add("email", "cadcsd")

	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("Got valid email when should have been invalid")
	}
}
