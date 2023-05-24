package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	"github.com/JakeStrang1/escapehatch/app"
	"github.com/JakeStrang1/escapehatch/email"
	api "github.com/JakeStrang1/escapehatch/http"
	"github.com/JakeStrang1/escapehatch/services/auth"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
)

type Suite struct {
	suite.Suite
	App   app.App
	Email string
}

/*******************************************
 *
 * testify/suite hooks
 * https://pkg.go.dev/github.com/stretchr/testify/suite
 *
 *******************************************/

func (s *Suite) SetupSuite() {
	// Set working directory to the backend folder
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(filepath.Join(dir, ".."))
	if err != nil {
		panic(err)
	}

	err = godotenv.Overload(".env") // Will overwrite existing env vars
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	s.Email = os.Getenv("TEST_EMAIL") // Used for generating seed data

	config := app.Config{
		MongoHost:         os.Getenv("MONGO_HOST"),
		MongoDatabaseName: os.Getenv("TEST_MONGO_DB_NAME"), // Different from local DB name to avoid overwriting
		CORSAllowOrigin:   os.Getenv("ORIGIN"),
		UseSendGrid:       "false", // Tests will never send real emails
	}
	s.App = app.NewApp(config)
}

func (s *Suite) TearDownSuite() {
	s.App.Close()
}

func (s *Suite) SetupTest() {
	_, _, db, err := mgm.DefaultConfigs()
	if err != nil {
		panic(err)
	}

	err = db.Drop(mgm.Ctx())
	if err != nil {
		panic(err)
	}
}

/*******************************************
 *
 * Custom suite methods
 *
 *******************************************/

type Request struct {
	*http.Request
}

type Response struct {
	Status   int
	Body     string
	Recorder *httptest.ResponseRecorder
}

func (r *Response) Debug() {
	fmt.Printf("Status: %d. Body: %+v\n", r.Status, r.Body)
}

func (r *Response) Unmarshal(s *Suite, obj interface{}) {
	err := json.Unmarshal([]byte(gjson.Get(r.Body, "data").String()), obj)
	s.Assert().NoError(err)
}

func (s *Suite) Get(path string, options ...func(*Request)) Response {
	w := httptest.NewRecorder()
	httpRequest, err := http.NewRequest("GET", path, nil)
	s.Assert().NoError(err)

	request := Request{Request: httpRequest}

	for _, option := range options {
		option(&request)
	}

	s.App.Router().ServeHTTP(w, request.Request)

	return Response{
		Status:   w.Code,
		Body:     w.Body.String(),
		Recorder: w,
	}
}

func (s *Suite) Post(path string, body interface{}, options ...func(*Request)) Response {
	w := httptest.NewRecorder()
	bodyBytes, err := json.Marshal(body)
	s.Assert().NoError(err)

	reader := bytes.NewReader(bodyBytes)
	httpRequest, err := http.NewRequest("POST", path, reader)
	s.Assert().NoError(err)

	request := Request{Request: httpRequest}

	for _, option := range options {
		option(&request)
	}

	s.App.Router().ServeHTTP(w, request.Request)

	return Response{
		Status:   w.Code,
		Body:     w.Body.String(),
		Recorder: w,
	}
}

func (s *Suite) Patch(path string, body interface{}, options ...func(*Request)) Response {
	w := httptest.NewRecorder()
	bodyBytes, err := json.Marshal(body)
	s.Assert().NoError(err)

	reader := bytes.NewReader(bodyBytes)
	httpRequest, err := http.NewRequest("PATCH", path, reader)
	s.Assert().NoError(err)

	request := Request{Request: httpRequest}

	for _, option := range options {
		option(&request)
	}

	s.App.Router().ServeHTTP(w, request.Request)

	return Response{
		Status:   w.Code,
		Body:     w.Body.String(),
		Recorder: w,
	}
}

func (s *Suite) Delete(path string, options ...func(*Request)) Response {
	w := httptest.NewRecorder()
	httpRequest, err := http.NewRequest("DELETE", path, nil)
	s.Assert().NoError(err)

	request := Request{Request: httpRequest}

	for _, option := range options {
		option(&request)
	}

	s.App.Router().ServeHTTP(w, request.Request)

	return Response{
		Status:   w.Code,
		Body:     w.Body.String(),
		Recorder: w,
	}
}

func WithCookie(cookie *http.Cookie) func(*Request) {
	return func(r *Request) {
		r.AddCookie(cookie)
	}
}

// NewSeedEmail takes a label like "user123" and the test email like "myemail@gmail.com" and returns "myemail+user123@gmail.com".
// Use this to generate seed accounts that you can log into.
func (s *Suite) NewSeedEmail(label string) string {
	ss := strings.Split(s.Email, "@")
	if len(ss) != 2 {
		s.Assert().Failf("Test email must be a valid email address: %s", s.Email)
	}

	return fmt.Sprintf("%s+%s@%s", ss[0], label, ss[1])
}

// CreateUser signs up a new user and returns the user as well as the cookie option to use in future requests.
func (s *Suite) CreateUser(emailAddress string) (*api.UserAPI, func(*Request)) {
	signUpBody := api.SignInBody{
		Email: lo.ToPtr(emailAddress),
	}
	response := s.Post("/auth/sign-up", signUpBody)
	s.Assert().Equal(200, response.Status)

	mailer := email.GetMailer().(*email.MockMailer)
	secret, err := auth.GetSecretFromEmail(mailer.SendParams.PlainContent)
	s.Assert().NoError(err)

	verifyBody := api.VerifyBody{
		Email:  lo.ToPtr(emailAddress),
		Secret: &secret,
	}
	response = s.Post("/auth/verify", verifyBody)
	s.Assert().Equal(200, response.Status)
	s.Assert().NotEmpty(response.Recorder.Result().Cookies()[0].Value)

	withCookie := WithCookie(response.Recorder.Result().Cookies()[0])
	response = s.Get("/users/me", withCookie)
	s.Assert().Equal(200, response.Status)

	user := api.UserAPI{}
	response.Unmarshal(s, &user)
	return &user, withCookie
}

func (s *Suite) DebugAndFail(response Response) {
	response.Debug()
	s.Assert().Fail("debug and fail")
}
