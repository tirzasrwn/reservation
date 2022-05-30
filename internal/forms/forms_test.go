package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)
	isValid := form.Valid()
	if !isValid {
		t.Error("Got invalid when should have been valid.")
	}
}

func TestFor(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fileds is missiong")
	}

	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/some-url", nil)
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Shows does not have required fields when it does.")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	has := form.Has("whatever")
	if has {
		t.Error("Field shows has field when it does not.")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("Show form does not have when it should have.")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/some-url", nil)
	form := New(r.PostForm)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("Form shows min length for none existent field.")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("Should have an error, but not have one.")
	}

	postedValue := url.Values{}
	postedValue.Add("a", "12345")
	form = New(postedValue)
	form.MinLength("a", 100)
	if form.Valid() {
		t.Error("Shows min length 100 met when data is shorter.")
	}

	postedValue = url.Values{}
	postedValue.Add("b", "12345")
	form = New(postedValue)
	form.MinLength("b", 3)
	if !form.Valid() {
		t.Error("Shows min length of 3 met when it is")
	}

	isError = form.Errors.Get("b")
	if isError != "" {
		t.Error("Should not have an error, but not got one.")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedValue := url.Values{}
	form := New(postedValue)
	form.IsEmail("c")
	if form.Valid() {
		t.Error("Form shows is email for none existent field.")
	}

	postedValue = url.Values{}
	postedValue.Add("d", "d")
	form = New(postedValue)
	form.IsEmail("d")
	if form.Valid() {
		t.Error("Shows valid email when it should not.")
	}

	postedValue = url.Values{}
	postedValue.Add("e", "arya@stark.com")
	form = New(postedValue)
	form.IsEmail("e")
	if !form.Valid() {
		t.Error("Shows unvalid email when it should valid.")
	}
}
