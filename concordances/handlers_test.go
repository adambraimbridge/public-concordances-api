package concordances

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server          *httptest.Server
	concordanceURL  string
	isFound         bool
	conceptIds      []string
	authorities     []string
	actualAuthority string
)

type mockConcordanceDriver struct{}

func (driver mockConcordanceDriver) ReadByConceptID(ids []string) (concordances Concordances, found bool, err error) {
	conceptIds = ids
	return Concordances{}, isFound, nil
}
func (driver mockConcordanceDriver) ReadByAuthority(authority string, ids []string) (concordances Concordances, found bool, err error) {
	authorities = ids
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
	assert.Len(authorities, 1)
	assert.Contains(authorities, "some-value")
}

func TestCanGetMultipleIdentifiersByAuthority(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	req, _ := http.NewRequest("GET", concordanceURL+"?authority=some-authority&identifierValue=some-value&identifierValue=some-value2", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.EqualValues(actualAuthority, "some-authority")
	assert.Len(authorities, 2)
	assert.Contains(authorities, "some-value")
	assert.Contains(authorities, "some-value2")
}

func TestReturnBadRequestGivenMoreThanOneAuthority(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	req, _ := http.NewRequest("GET", concordanceURL+"?authority=some-authority&identifierValue=some-value&authority=some-authority-yet-again", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(400, res.StatusCode)
}

func TestCanGetOneConcept(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=bob", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.Len(conceptIds, 1)
	assert.Contains(conceptIds, "bob")
}

func TestCanGetMultipleConcepts(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=bob&conceptId=carlos", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.Len(conceptIds, 2)
	assert.Contains(conceptIds, "carlos")
}

func TestCanParseConceptURI(t *testing.T) {
	assert := assert.New(t)
	isFound = true
	req, _ := http.NewRequest("GET", concordanceURL+"?conceptId=http://api.ft.com/things/8138ca3f-b80d-3ef8-ad59-6a9b6ea5f15e", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.EqualValues(200, res.StatusCode)
	assert.Len(conceptIds, 1)
	assert.Contains(conceptIds, "8138ca3f-b80d-3ef8-ad59-6a9b6ea5f15e")
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
