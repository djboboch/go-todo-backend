package models

import (
	"errors"
	"testing"
)

func TestCreatePost(t *testing.T) {

	t.Run("Missing Header", func(t *testing.T) {
		gotPost, gotError := CreatePost("", "Todo body")

		wantPost := Post{}
		wantError := errors.New(ErrorMissingHeader)

		if gotPost != wantPost {
			t.Errorf("got '%+v' post body containing data when '%+v' wanted empty", gotPost, wantPost)
		}

		if gotError != wantError {
			t.Errorf("got '%v' error message wanted '%v' error message", gotError, wantError)
		}
	})

	t.Run("Correct Format", func(t *testing.T) {

		const (
			HeaderText = "Shopping List"
			BodyText   = "Remember to buy: apples, banana"
		)

		got, _ := CreatePost(HeaderText, BodyText)

		want := Post{HeaderText, BodyText}

		if got != want {
			t.Errorf("got %+v want %+v", got, want)
		}
	})
}
