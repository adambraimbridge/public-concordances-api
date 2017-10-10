package concordances

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"reflect"

	"sort"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/concepts-rw-neo4j/concepts"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/Financial-Times/organisations-rw-neo4j/organisations"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

var concordedBrandSmartlogic = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/b20801ac-5a76-43cf-b816-8c3b2f7133ad",
		APIURL: "http://api.ft.com/brands/b20801ac-5a76-43cf-b816-8c3b2f7133ad"},
	Identifier{
		Authority:       "http://api.ft.com/system/SMARTLOGIC",
		IdentifierValue: "b20801ac-5a76-43cf-b816-8c3b2f7133ad"},
}

var concordedBrandSmartlogicUPP = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/b20801ac-5a76-43cf-b816-8c3b2f7133ad",
		APIURL: "http://api.ft.com/brands/b20801ac-5a76-43cf-b816-8c3b2f7133ad"},
	Identifier{
		Authority:       "http://api.ft.com/system/UPP",
		IdentifierValue: "b20801ac-5a76-43cf-b816-8c3b2f7133ad"},
}

var concordedBrandTME = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/b20801ac-5a76-43cf-b816-8c3b2f7133ad",
		APIURL: "http://api.ft.com/brands/b20801ac-5a76-43cf-b816-8c3b2f7133ad"},
	Identifier{
		Authority:       "http://api.ft.com/system/FT-TME",
		IdentifierValue: "VGhlIFJvbWFu-QnJhbmRz"},
}

var concordedBrandTMEUPP = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/b20801ac-5a76-43cf-b816-8c3b2f7133ad",
		APIURL: "http://api.ft.com/brands/b20801ac-5a76-43cf-b816-8c3b2f7133ad"},
	Identifier{
		Authority:       "http://api.ft.com/system/UPP",
		IdentifierValue: "70f4732b-7f7d-30a1-9c29-0cceec23760e"},
}

var mainOrganisationLEI = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/21a5cc0-d326-4e62-b84a-d840c2209fee",
		APIURL: "http://api.ft.com/organisations/f21a5cc0-d326-4e62-b84a-d840c2209fee"},
	Identifier{
		Authority:       "http://api.ft.com/system/LEI",
		IdentifierValue: "7ZW8QJWVPR4P1J1KQY45"},
}

var childOrganisationFactset = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/f21a5cc0-d326-4e62-b84a-d840c2209fee",
		APIURL: "http://api.ft.com/organisations/f21a5cc0-d326-4e62-b84a-d840c2209fee"},
	Identifier{
		Authority:       "http://api.ft.com/system/FACTSET",
		IdentifierValue: "003JLG-E"},
}

var childOrganisationFactsetUPP = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/f21a5cc0-d326-4e62-b84a-d840c2209fee",
		APIURL: "http://api.ft.com/organisations/f21a5cc0-d326-4e62-b84a-d840c2209fee"},
	Identifier{
		Authority:       "http://api.ft.com/system/UPP",
		IdentifierValue: "f21a5cc0-d326-4e62-b84a-d840c2209fee"},
}

var childOrganisationLEI = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/f21a5cc0-d326-4e62-b84a-d840c2209fee",
		APIURL: "http://api.ft.com/organisations/f21a5cc0-d326-4e62-b84a-d840c2209fee"},
	Identifier{
		Authority:       "http://api.ft.com/system/LEI",
		IdentifierValue: "7ZW8QJWVPR4P1J1KQY45"},
}

var mainOrganisationTME = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/3e844449-b27f-40d4-b696-2ce9b6137133",
		APIURL: "http://api.ft.com/organisations/3e844449-b27f-40d4-b696-2ce9b6137133"},
	Identifier{
		Authority:       "http://api.ft.com/system/FT-TME",
		IdentifierValue: "TnN0ZWluX09OX0ZvcnR1bmVDb21wYW55X05XUw==-T04="},
}

var mainOrganisationTMEUPP = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/3e844449-b27f-40d4-b696-2ce9b6137133",
		APIURL: "http://api.ft.com/organisations/3e844449-b27f-40d4-b696-2ce9b6137133"},
	Identifier{
		Authority:       "http://api.ft.com/system/UPP",
		IdentifierValue: "TnN0ZWluX09OX0ZvcnR1bmVDb21wYW55X05XUw==-T04="},
}

var unconcordedBrandTME = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/ad56856a-7d38-48e2-a131-7d104f17e8f6",
		APIURL: "http://api.ft.com/brands/ad56856a-7d38-48e2-a131-7d104f17e8f6"},
	Identifier{
		Authority:       "http://api.ft.com/system/FT-TME",
		IdentifierValue: "UGFydHkgcGVvcGxl-QnJhbmRz"},
}

var unconcordedBrandTMEUPP = Concordance{
	Concept{
		ID:     "http://api.ft.com/things/ad56856a-7d38-48e2-a131-7d104f17e8f6",
		APIURL: "http://api.ft.com/brands/ad56856a-7d38-48e2-a131-7d104f17e8f6"},
	Identifier{
		Authority:       "http://api.ft.com/system/UPP",
		IdentifierValue: "ad56856a-7d38-48e2-a131-7d104f17e8f6"},
}

func TestNeoReadByConceptID_NewModel_Unconcorded(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)

	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())

	writeGenericConceptJSONToService(conceptRW, "./fixtures/Brand-Unconcorded-ad56856a-7d38-48e2-a131-7d104f17e8f6.json", assert)
	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByConceptID([]string{"ad56856a-7d38-48e2-a131-7d104f17e8f6"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(2, len(conc.Concordance))

	readConceptAndCompare(t, Concordances{[]Concordance{unconcordedBrandTME, unconcordedBrandTMEUPP}}, conc, "TestNeoReadByConceptID_NewModel_Unconcorded")
}

func TestNeoReadByConceptID_NewModel_Concorded(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())

	writeGenericConceptJSONToService(conceptRW, "./fixtures/Brand-Concorded-b20801ac-5a76-43cf-b816-8c3b2f7133ad.json", assert)
	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByConceptID([]string{"b20801ac-5a76-43cf-b816-8c3b2f7133ad"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(4, len(conc.Concordance))

	readConceptAndCompare(t, Concordances{[]Concordance{concordedBrandSmartlogic, concordedBrandSmartlogicUPP, concordedBrandTME, concordedBrandTMEUPP}}, conc, "TestNeoReadByConceptID_NewModel_Concorded")
}

func TestNeoReadByConceptID_NewModel_And_OldModel(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())
	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())

	writeJSONToService(organisationRW, "./fixtures/Organisation-Child-f21a5cc0-d326-4e62-b84a-d840c2209fee.json", assert)
	writeGenericConceptJSONToService(conceptRW, "./fixtures/Brand-Concorded-b20801ac-5a76-43cf-b816-8c3b2f7133ad.json", assert)

	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByConceptID([]string{"b20801ac-5a76-43cf-b816-8c3b2f7133ad", "f21a5cc0-d326-4e62-b84a-d840c2209fee"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(7, len(conc.Concordance))

	sliceConcordances := []Concordance{childOrganisationFactset, childOrganisationFactsetUPP, childOrganisationLEI, concordedBrandTME, concordedBrandTMEUPP, concordedBrandSmartlogic, concordedBrandSmartlogicUPP}
	readConceptAndCompare(t, Concordances{sliceConcordances}, conc, "TestNeoReadByConceptID_NewModel_And_OldModel")
}

func TestNeoReadByAuthority_NewModel_Unconcorded(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())

	writeGenericConceptJSONToService(conceptRW, "./fixtures/Brand-Unconcorded-ad56856a-7d38-48e2-a131-7d104f17e8f6.json", assert)
	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByAuthority("http://api.ft.com/system/FT-TME", []string{"UGFydHkgcGVvcGxl-QnJhbmRz"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(1, len(conc.Concordance))

	sliceConcordances := []Concordance{unconcordedBrandTME}
	readConceptAndCompare(t, Concordances{sliceConcordances}, conc, "TestNeoReadByAuthority_NewModel_Unconcorded")
}

func TestNeoReadByAuthority_NewModel_Concorded(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())

	writeGenericConceptJSONToService(conceptRW, "./fixtures/Brand-Concorded-b20801ac-5a76-43cf-b816-8c3b2f7133ad.json", assert)
	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByAuthority("http://api.ft.com/system/SMARTLOGIC", []string{"b20801ac-5a76-43cf-b816-8c3b2f7133ad"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(1, len(conc.Concordance))

	readConceptAndCompare(t, Concordances{[]Concordance{concordedBrandSmartlogic}}, conc, "TestNeoReadByAuthority_NewModel_Concorded")
}

func TestNeoReadByAuthority_NewModel_And_OldModel(t *testing.T) {
	assert := assert.New(t)
	db := getDatabaseConnection(t, assert)
	conceptRW := concepts.NewConceptService(db)
	assert.NoError(conceptRW.Initialise())
	organisationRW := organisations.NewCypherOrganisationService(db)
	assert.NoError(organisationRW.Initialise())

	writeJSONToService(organisationRW, "./fixtures/Organisation-Main-3e844449-b27f-40d4-b696-2ce9b6137133.json", assert)
	writeGenericConceptJSONToService(conceptRW, "./fixtures/Brand-Concorded-b20801ac-5a76-43cf-b816-8c3b2f7133ad.json", assert)

	defer cleanUp(assert, db)

	undertest := NewCypherDriver(db, "prod")
	conc, found, err := undertest.ReadByAuthority("http://api.ft.com/system/FT-TME", []string{"TnN0ZWluX09OX0ZvcnR1bmVDb21wYW55X05XUw==-T04=", "VGhlIFJvbWFu-QnJhbmRz"})
	assert.NoError(err)
	assert.True(found)
	assert.Equal(2, len(conc.Concordance))

	readConceptAndCompare(t, Concordances{[]Concordance{concordedBrandTME, mainOrganisationTME}}, conc, "TestNeoReadByAuthority_NewModel_And_OldModel")
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

	readConceptAndCompare(t, Concordances{[]Concordance{childOrganisationFactset, childOrganisationFactsetUPP, childOrganisationLEI}}, cs, "TestNeoReadByConceptIDToConcordancesMandatoryFields")
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

	readConceptAndCompare(t, Concordances{[]Concordance{childOrganisationFactset}}, cs, "TestNeoReadByAuthorityToConcordancesMandatoryFields")
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

	readConceptAndCompare(t, Concordances{[]Concordance{childOrganisationFactset}}, cs, "TestNeoReadByAuthorityOnlyOneConcordancePerIdentifierValue")
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
	assert.Equal(3, len(cs.Concordance))

	readConceptAndCompare(t, Concordances{[]Concordance{childOrganisationFactset, childOrganisationFactsetUPP, childOrganisationLEI}}, cs, "TestNeoReadByConceptIdReturnMultipleConcordancesForMultipleIdentifiers")
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
}

func readConceptAndCompare(t *testing.T, expected Concordances, actual Concordances, testName string) {

	fullConcordanceSort(expected.Concordance)
	fullConcordanceSort(actual.Concordance)

	assert.True(t, reflect.DeepEqual(expected, actual), fmt.Sprintf("Actual aggregated concept differs from expected: Test: %v \n Expected: %v \n Actual: %v", testName, expected, actual))
}

func fullConcordanceSort(concordanceList []Concordance) {
	sort.SliceStable(concordanceList, func(i, j int) bool {
		return concordanceList[i].Concept.ID < concordanceList[j].Concept.ID
	})
	sort.SliceStable(concordanceList, func(i, j int) bool {
		return concordanceList[i].Identifier.Authority < concordanceList[j].Identifier.Authority
	})
	sort.SliceStable(concordanceList, func(i, j int) bool {
		return concordanceList[i].Identifier.IdentifierValue < concordanceList[j].Identifier.IdentifierValue
	})
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

func writeGenericConceptJSONToService(service concepts.ConceptService, pathToJSONFile string, assert *assert.Assertions) {
	f, err := os.Open(pathToJSONFile)
	assert.NoError(err)
	dec := json.NewDecoder(f)
	inst, _, errr := service.DecodeJSON(dec)
	assert.NoError(errr)
	_, errrr := service.Write(inst, "test_transaction_id")
	assert.NoError(errrr)
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
