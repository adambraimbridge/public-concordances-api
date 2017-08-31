package concordances

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/concepts-rw-neo4j/concepts"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/Financial-Times/organisations-rw-neo4j/organisations"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

func TestNeoReadByConceptID_NewModel_Unconcorded(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)

	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())

	writeJSONToService(conceptRW, "./fixtures/Brand-Unconcorded-ad56856a-7d38-48e2-a131-7d104f17e8f6.json", assert)
	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByConceptID([]string{"ad56856a-7d38-48e2-a131-7d104f17e8f6"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(1, len(conc.Concordance))
}

func TestNeoReadByConceptID_NewModel_Concorded(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())

	writeJSONToService(conceptRW, "./fixtures/Brand-Concorded-b20801ac-5a76-43cf-b816-8c3b2f7133ad.json", assert)
	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByConceptID([]string{"b20801ac-5a76-43cf-b816-8c3b2f7133ad"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(2, len(conc.Concordance))
}

func TestNeoReadByConceptID_NewModel_And_OldModel(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())
	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())

	writeJSONToService(organisationRW, "./fixtures/Organisation-Child-f21a5cc0-d326-4e62-b84a-d840c2209fee.json", assert)
	writeJSONToService(conceptRW, "./fixtures/Brand-Concorded-b20801ac-5a76-43cf-b816-8c3b2f7133ad.json", assert)

	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByConceptID([]string{"b20801ac-5a76-43cf-b816-8c3b2f7133ad", "f21a5cc0-d326-4e62-b84a-d840c2209fee"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(4, len(conc.Concordance))
}

func TestNeoReadByAuthority_NewModel_Unconcorded(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())

	writeJSONToService(conceptRW, "./fixtures/Brand-Unconcorded-ad56856a-7d38-48e2-a131-7d104f17e8f6.json", assert)
	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByAuthority("http://api.ft.com/system/FT-TME", []string{"UGFydHkgcGVvcGxl-QnJhbmRz"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(1, len(conc.Concordance))
}

func TestNeoReadByAuthority_NewModel_Concorded(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())

	writeJSONToService(conceptRW, "./fixtures/Brand-Concorded-b20801ac-5a76-43cf-b816-8c3b2f7133ad.json", assert)
	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByAuthority("http://api.ft.com/system/SMARTLOGIC", []string{"b20801ac-5a76-43cf-b816-8c3b2f7133ad"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(1, len(conc.Concordance))
}

func TestNeoReadByAuthority_NewModel_And_OldModel(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())
	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())

	writeJSONToService(organisationRW, "./fixtures/Organisation-Main-3e844449-b27f-40d4-b696-2ce9b6137133.json", assert)
	writeJSONToService(conceptRW, "./fixtures/Brand-Concorded-b20801ac-5a76-43cf-b816-8c3b2f7133ad.json", assert)

	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByAuthority("http://api.ft.com/system/FT-TME", []string{"TnN0ZWluX09OX0ZvcnR1bmVDb21wYW55X05XUw==-T04=", "VGhlIFJvbWFu-QnJhbmRz"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(2, len(conc.Concordance))
}

func TestNeoReadByConceptIDToConcordancesMandatoryFields(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())

	writeJSONToService(organisationRW, "./fixtures/Organisation-Child-f21a5cc0-d326-4e62-b84a-d840c2209fee.json", assert)

	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	cs, found, err := undertest.ReadByConceptID([]string{"f21a5cc0-d326-4e62-b84a-d840c2209fee"})
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(cs.Concordance)
	cleanUpParentOrgAndUppIdentifier(db, t, assert)
}

func TestNeoReadByAuthorityToConcordancesMandatoryFields(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)

	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())

	writeJSONToService(organisationRW, "./fixtures/Organisation-Child-f21a5cc0-d326-4e62-b84a-d840c2209fee.json", assert)

	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	cs, found, err := undertest.ReadByAuthority("http://api.ft.com/system/FACTSET", []string{"003JLG-E"})
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(cs.Concordance)

}

func TestNeoReadByAuthorityOnlyOneConcordancePerIdentifierValue(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)

	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())

	writeJSONToService(organisationRW, "./fixtures/Organisation-Child-f21a5cc0-d326-4e62-b84a-d840c2209fee.json", assert)

	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	cs, found, err := undertest.ReadByAuthority("http://api.ft.com/system/FACTSET", []string{"003JLG-E"})
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(cs.Concordance)
	assert.Equal(len(cs.Concordance), 1)
}

func TestNeoReadByConceptIdReturnMultipleConcordancesForMultipleIdentifiers(t *testing.T) {

	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)

	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())

	writeJSONToService(organisationRW, "./fixtures/Organisation-Child-f21a5cc0-d326-4e62-b84a-d840c2209fee.json", assert)

	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	cs, found, err := undertest.ReadByConceptID([]string{"f21a5cc0-d326-4e62-b84a-d840c2209fee"})
	assert.NoError(err)
	assert.True(found)
	assert.NotEmpty(cs.Concordance)
	assert.Equal(len(cs.Concordance), 2)

}

func TestNeoReadByAuthorityEmptyConcordancesWhenUnsupportedAuthority(t *testing.T) {

	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)

	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())

	writeJSONToService(organisationRW, "./fixtures/Organisation-Child-f21a5cc0-d326-4e62-b84a-d840c2209fee.json", assert)

	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	cs, found, err := undertest.ReadByAuthority("http://api.ft.com/system/UnsupportedAuthority", []string{"DANMUR-1"})
	assert.NoError(err)
	assert.False(found)
	assert.Empty(cs.Concordance)
	cleanUpParentOrgAndUppIdentifier(db, t, assert)
}

func getDatabaseConnection(t *testing.T, assert *assert.Assertions) neoutils.NeoConnection {
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	conf := neoutils.DefaultConnectionConfig()
	conf.Transactional = false
	db, err := neoutils.Connect(url, conf)
	assert.NoError(err, "Failed to connect to Neo4j")
	return db
}

func writeJSONToService(service baseftrwapp.Service, pathToJSONFile string, assert *assert.Assertions) {
	f, err := os.Open(pathToJSONFile)
	assert.NoError(err)
	dec := json.NewDecoder(f)
	inst, _, errr := service.DecodeJSON(dec)
	assert.NoError(errr)
	errrr := service.Write(inst, "test_transaction_id")
	assert.NoError(errrr)
}

func cleanUp(assert *assert.Assertions, db neoutils.NeoConnection) {

	qs := []*neoism.CypherQuery{
		{
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%s'})--(i:Identifier) DETACH DELETE i, a", "3e844449-b27f-40d4-b696-2ce9b6137133"),
		}, {
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%s'})--(i:Identifier) DETACH DELETE i, a", "f21a5cc0-d326-4e62-b84a-d840c2209fee"),
		}, {
			Statement: fmt.Sprintf("MATCH (a:Thing {uuid: '%s'})--(i:Identifier) DETACH DELETE i, a", "f9694ba7-eab0-4ce0-8e01-ff64bccb813c"),
		}, {
			Statement: fmt.Sprintf("MATCH (t:Thing {uuid: '%v'})--(i:Identifier) OPTIONAL MATCH (t)-[:EQUIVALENT_TO]-(e:Thing) DETACH DELETE t, i", "70f4732b-7f7d-30a1-9c29-0cceec23760e"),
		}, {
			Statement: fmt.Sprintf("MATCH (t:Thing {uuid: '%v'})--(i:Identifier) OPTIONAL MATCH (t)-[:EQUIVALENT_TO]-(e:Thing) DETACH DELETE t, e, i", "b20801ac-5a76-43cf-b816-8c3b2f7133ad"),
		}, {
			Statement: fmt.Sprintf("MATCH (t:Thing {uuid: '%v'})--(i:Identifier) OPTIONAL MATCH (t)-[:EQUIVALENT_TO]-(e:Thing) DETACH DELETE t, i", "ad56856a-7d38-48e2-a131-7d104f17e8f6"),
		}, {
			Statement: fmt.Sprintf("MATCH (t:Thing {uuid: '%v'})--(i:Identifier) OPTIONAL MATCH (t)-[:EQUIVALENT_TO]-(e:Thing) DETACH DELETE t, i", "dbb0bdae-1f0c-11e4-b0cb-b2227cce2b54"),
		}, {
			Statement: fmt.Sprintf("MATCH (t:Thing {prefUUID: '%v'}) DETACH DELETE t", "ad56856a-7d38-48e2-a131-7d104f17e8f6"),
		},
	}

	err := db.CypherBatch(qs)
	assert.NoError(err)

}

func cleanUpParentOrgAndUppIdentifier(db neoutils.NeoConnection, t *testing.T, assert *assert.Assertions) {
	qs := []*neoism.CypherQuery{
		{
			//deletes parent 'org' which only has type Thing
			Statement: fmt.Sprintf("MATCH (j:Thing {uuid: '%v'}) DETACH DELETE j", "3e844449-b27f-40d4-b696-2ce9b6137133"),
		},
		{
			//deletes upp identifier for the above parent 'org'
			Statement: fmt.Sprintf("MATCH (k:Identifier {value: '%v'}) DETACH DELETE k", "3e844449-b27f-40d4-b696-2ce9b6137133"),
		},
	}

	err := db.CypherBatch(qs)
	assert.NoError(err)
}
