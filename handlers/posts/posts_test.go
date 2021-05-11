package posts

import (
	"encoding/json"
	"github.com/djboboch/go-todo/internal/http/responses"
	"github.com/djboboch/go-todo/models"
	"github.com/segmentio/ksuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPostModel struct{}

var apiPrefix = "/api/v1"

func (m mockPostModel) All() ([]models.Post, error) {
	var posts []models.Post

	posts = append(posts, models.Post{
		Id:             ksuid.New().String(),
		Content:        "todo one",
		IsItemFinished: false,
	})

	posts = append(posts, models.Post{
		Id:             ksuid.New().String(),
		Content:        "todo two",
		IsItemFinished: true,
	})

	return posts, nil
}

func (m mockPostModel) Create(content string) (*models.Post, error) {
	return &models.Post{
		Id:             ksuid.New().String(),
		Content:        content,
		IsItemFinished: false,
	}, nil
}

func (m mockPostModel) Update(post models.Post) error {
	return nil
}

func (m mockPostModel) Delete(id string) error {
	return nil
}

func TestEnv_GetPosts(t *testing.T) {

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", apiPrefix+"/todo", nil)
	env := Env{Posts: mockPostModel{}}

	env.GetPosts().ServeHTTP(res, req)

	var serverResponse responses.ServerResponse

	json.Unmarshal(res.Body.Bytes(), &serverResponse)

	t.Run("Request Content-Type", func(t *testing.T) {
		if res.Code != http.StatusOK {
			t.Errorf("Status Code should be 200")
		}
	})

	t.Run("Request Response Status", func(t *testing.T) {
		expectedStatus := responses.SuccessResponseStatus

		if serverResponse.Status != expectedStatus {
			t.Errorf("Exptected status: %v\n...obtained = %v", expectedStatus, serverResponse.Status)
		}
	})
}

func TestEnv_CreatePost(t *testing.T) {

}
