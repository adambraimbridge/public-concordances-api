package concordances

import (
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestNeoReadStructToPersonMandatoryFields checks that madatory fields are set even if they are empty or nil / null
func TestNeoReadStructToConcordancesMandatoryFields(t *testing.T) {

	assert := assert.New(t)
	db := getDatabaseConnectionAndCheckClean(t, assert)
	//	batchRunner := neoutils.NewBatchCypherRunner(neoutils.StringerDb{db}, 1)

	undertest := NewCypherDriver(db, "prod")
	cs, found, err := undertest.ReadByConceptId([]string{"3e844449-b27f-40d4-b696-2ce9b6137133"})
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(cs.Concordance)
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

}
