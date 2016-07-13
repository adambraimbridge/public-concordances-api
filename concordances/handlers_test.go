package concordances

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server          *httptest.Server
	concordanceURL  string
	isFound         bool
	conceptId       string
	authorityValue  string
	actualAuthority string
	isCanonical     bool
)

const (
	alternativeId = "alternativeId"
)

type mockConcordanceDriver struct{}

func (driver mockConcordanceDriver) ReadByConceptID(id string) (concordances Concordances, found bool, err error) {
	if isCanonical {
		conceptId = id
	} else {
		conceptId = alternativeId
	}
	return Concordances{[]Concordance{Concordance{Concept: Concept{ID: conceptId}}}}, isFound, nil
}
func (driver mockConcordanceDriver) ReadByAuthority(authority string, id string) (concordances Concordances, found bool, err error) {
	authorityValue = id
	actualAuthority = authority
	return Concordances{}, isFound, nil
}

func (driver mockConcordanceDriver) CheckConnectivity() error {
	return nil
}

func init() {
	ConcordanceDriver = mockConcordanceDriver{}
	r := mux.NewRouter()
	r.HandleFunc("/concordances", GetConcordances).Methods("GET")
	server = httptest.NewServer(r)
	concordanceURL = fmt.Sprintf("%s/concordances", server.URL) //Grab the address for the API endpoint
	isFound = true
	isCanonical = true
}

func TestCanGetOneAuthority(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	isCanonical = true
	req, _ := http.NewRequest("GET", concordanceURL+"?authority=some-authority&identifierValue=some-value", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.EqualValues(actualAuthority, "some-authority")
	assert.Equal("some-value", authorityValue)
}

func TestCanNotGetMultipleIdentifiersByAuthority(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	isCanonical = true
	req, _ := http.NewRequest("GET", concordanceURL+"?authority=some-authority&identifierValue=some-value&identifierValue=some-value2", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)

	assert.EqualValues(400, res.StatusCode)
	msg, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)
	assert.Contains(string(msg), multipleAuthorityValuesNotSupported)
}

func TestReturnBadRequestGivenMoreThanOneAuthority(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	isCanonical = true
	req, _ := http.NewRequest("GET", concordanceURL+"?authority=some-authority&identifierValue=some-value&authority=some-authority-yet-again", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(400, res.StatusCode)
	msg, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)
	assert.Contains(string(msg), multipleAuthoritiesNotPermitted)
}

func TestCanGetOneConcept(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	isCanonical = true
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=bob", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.Equal("bob", conceptId)
}

func TestCanNotGetMultipleConcepts(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	isCanonical = true
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=bob&conceptId=carlos", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(400, res.StatusCode)
	msg, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)
	assert.Contains(string(msg), multipleConceptIDsNotSupported)
}

func TestCanParseConceptURI(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	isCanonical = true
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=http://api.ft.com/things/8138ca3f-b80d-3ef8-ad59-6a9b6ea5f15e", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.Equal("8138ca3f-b80d-3ef8-ad59-6a9b6ea5f15e", conceptId)
}

func TestCanNotRequestAuthorityAndConceptId(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	isCanonical = true
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=bob&authority=high-and-mighty", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(400, res.StatusCode)
}

func TestCanNotRequestWithoutAuthorityOrConceptId(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	isCanonical = true
	req, _ := http.NewRequest("GET", concordanceURL+"?randomRequestParam=bob", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(400, res.StatusCode)
}

func noRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("Don't redirect!")
}

func TestRedirectHappensOnFoundForAlternateNode(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	isCanonical = false
	request, _ := http.NewRequest("GET", concordanceURL+"?conceptId=bob", nil)
	cl := &http.Client{
		CheckRedirect: noRedirect,
	}
	result, err := cl.Do(request)

	assert.Contains(err.Error(), "Don't redirect!")
	assert.EqualValues(301, result.StatusCode)
	assert.Equal("/concordances?conceptId="+alternativeId, result.Header.Get("Location"))
	assert.Equal("application/json; charset=UTF-8", result.Header.Get("Content-Type"))
}
