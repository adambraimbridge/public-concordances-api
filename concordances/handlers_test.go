package concordances

import (
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
)

type mockConcordanceDriver struct{}

func (driver mockConcordanceDriver) ReadByConceptID(id string) (concordances Concordances, found bool, err error) {
	conceptId = id
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
}

func TestCanGetOneAuthority(t *testing.T) {
	assert := assert.New(t)
	isFound = true
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
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=bob", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.Equal("bob", conceptId)
}

func TestCanNotGetMultipleConcepts(t *testing.T) {
	assert := assert.New(t)
	isFound = true
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
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=http://api.ft.com/things/8138ca3f-b80d-3ef8-ad59-6a9b6ea5f15e", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.Equal("8138ca3f-b80d-3ef8-ad59-6a9b6ea5f15e", conceptId)
}

func TestCanNotRequestAuthorityAndConceptId(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=bob&authority=high-and-mighty", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(400, res.StatusCode)
}

func TestCanNotRequestWithoutAuthorityOrConceptId(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	req, _ := http.NewRequest("GET", concordanceURL+"?randomRequestParam=bob", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(400, res.StatusCode)
}
