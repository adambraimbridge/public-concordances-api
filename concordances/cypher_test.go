package concordances

import (
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/Financial-Times/organisations-rw-neo4j/organisations"
	"github.com/Financial-Times/people-rw-neo4j/people"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestNeoReadStructToPersonMandatoryFields checks that madatory fields are set even if they are empty or nil / null
func TestNeoReadByConceptIDToConcordancesMandatoryFields(t *testing.T) {

	assert := assert.New(t)
	db := getDatabaseConnectionAndCheckClean(t, assert)
	batchRunner := neoutils.NewBatchCypherRunner(neoutils.StringerDb{db}, 1)

	peopleRW, organisationRW := getServices(t, assert, db, &batchRunner)
	writeJSONToService(peopleRW, "./fixtures/Person-Dan_Murphy-868c3c17-611c-4943-9499-600ccded71f3.json", assert)
	writeJSONToService(organisationRW, "./fixtures/Organisation-Child-f21a5cc0-d326-4e62-b84a-d840c2209fee.json", assert)

	undertest := NewCypherDriver(db, "prod")
	cs, found, err := undertest.ReadByConceptID([]string{"868c3c17-611c-4943-9499-600ccded71f3"})
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(cs.Concordance)
	fmt.Printf("RESULTS:%s\n", cs)
}

func TestNeoReadByAuthorityToConcordancesMandatoryFields(t *testing.T) {

	assert := assert.New(t)
	db := getDatabaseConnectionAndCheckClean(t, assert)
	batchRunner := neoutils.NewBatchCypherRunner(neoutils.StringerDb{db}, 1)

	peopleRW, organisationRW := getServices(t, assert, db, &batchRunner)
	writeJSONToService(peopleRW, "./fixtures/Person-Dan_Murphy-868c3c17-611c-4943-9499-600ccded71f3.json", assert)
	writeJSONToService(organisationRW, "./fixtures/Organisation-Child-f21a5cc0-d326-4e62-b84a-d840c2209fee.json", assert)

	undertest := NewCypherDriver(db, "prod")
	cs, found, err := undertest.ReadByAuthority("http://api.ft.com/system/FACTSET-PPL", []string{"DANMUR-1"})
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(cs.Concordance)
	fmt.Printf("RESULTS:%s\n", cs)
}

func getServices(t *testing.T, assert *assert.Assertions, db *neoism.Database, batchRunner *neoutils.CypherRunner) (baseftrwapp.Service, baseftrwapp.Service) {
	peopleRW := people.NewCypherPeopleService(*batchRunner, db)
	assert.NoError(peopleRW.Initialise())
	organisationRW := organisations.NewCypherOrganisationService(*batchRunner, db)
	assert.NoError(organisationRW.Initialise())
	return peopleRW, organisationRW
}

func getDatabaseConnectionAndCheckClean(t *testing.T, assert *assert.Assertions) *neoism.Database {
	db := getDatabaseConnection(t, assert)
	cleanDB(db, t, assert)
	//	checkDbClean(db, t)
	return db
}

func getDatabaseConnection(t *testing.T, assert *assert.Assertions) *neoism.Database {
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	db, err := neoism.Connect(url)
	assert.NoError(err, "Failed to connect to Neo4j")
	return db
}

func cleanDB(db *neoism.Database, t *testing.T, assert *assert.Assertions) {
	uuids := []string{
		"f21a5cc0-d326-4e62-b84a-d840c2209fee",
		"868c3c17-611c-4943-9499-600ccded71f3",
	}

	qs := make([]*neoism.CypherQuery, len(uuids))
	for i, uuid := range uuids {
		qs[i] = &neoism.CypherQuery{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%s'}) DETACH DELETE a", uuid)}
	}
	err := db.CypherBatch(qs)
	assert.NoError(err)
}

func writeJSONToService(service baseftrwapp.Service, pathToJSONFile string, assert *assert.Assertions) {
	f, err := os.Open(pathToJSONFile)
	assert.NoError(err)
	dec := json.NewDecoder(f)
	inst, _, errr := service.DecodeJSON(dec)
	assert.NoError(errr)
	errrr := service.Write(inst)
	assert.NoError(errrr)
}
func deleteAllViaService(assert *assert.Assertions, peopleRW baseftrwapp.Service, organisationRW baseftrwapp.Service) {
	peopleRW.Delete("868c3c17-611c-4943-9499-600ccded71f3")
	organisationRW.Delete("f21a5cc0-d326-4e62-b84a-d840c2209fee")
}
