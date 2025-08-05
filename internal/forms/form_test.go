package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	form := getNewForm()
	if form.Valid() != true {
		t.Error("error found form is invalid")
	}
}

func getNewForm() *Form {
	data := url.Values{}
	data.Add("a", "aa")
	data.Add("b", "bb")
	data.Add("c", "")
	form := New(data)
	return form
}

func TestForm_Required(t *testing.T) {
	data := url.Values{}
	data.Add("a", "aa")
	data.Add("b", "bb")
	data.Add("c", "")
	form := New(data)
	fields := []string{"a", "b", "c"}
	form.Required(fields...)
	for _, s := range fields {
		if data.Get(s) == "" {
			if form.Errors.Get(s) != "This field cannot be blank" {
				t.Error("errors not getting added to form")
			}
		} else {
			if form.Errors.Get(s) != "" {
				t.Error("not expecting errors")
			}
		}
	}
}

func TestForm_MinLength(t *testing.T) {
	form := getNewForm()
	r := httptest.NewRequest("POST", "/some-url", nil)
	r.Form = getNewForm().Values
	form.MinLength("a", 2, r)
	if len(form.Errors) != 0 {
		t.Error("error in checking min length")
	}
	form.MinLength("c", 2, r)
	if form.Errors.Get("c") == "This field must be atleast 2 long" {
		t.Error("errors in checking min length")
	}
}

func TestForm_IsEmail(t *testing.T) {
	form := getNewForm()
	form.Add("email", "fijfi@email.com")
	form.IsEmail("email")
	if len(form.Errors) != 0 {
		t.Error("error validating email")
	}
	form.Set("email", "jfdofj")
	form.IsEmail("email")
	if form.Errors.Get("email") != "Invalid email address" {
		t.Error("error validating email")
	}
}

func TestForm_Has(t *testing.T) {
	form := getNewForm()
	r, _ := http.NewRequest("GET", "/some-url", nil)
	r.Form = form.Values
	if !form.Has("a", r) {
		t.Error("error checking field value in form")
	}
	if !form.Has("c", r) {
		if form.Errors.Get("c") != "this field cannot be blank" {
			t.Error("error checking empty field")
		}
	}
}
