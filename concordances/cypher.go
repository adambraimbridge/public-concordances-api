package concordances

import (
	"fmt"

	"github.com/Financial-Times/neo-model-utils-go/mapper"
	log "github.com/Sirupsen/logrus"
	"github.com/jmcvetta/neoism"
)

// Driver interface
type Driver interface {
	ReadByConceptId(ids []string) (concordances Concordances, found bool, err error)
	ReadByAuthority(authority string, ids []string) (concordances Concordances, found bool, err error)
	CheckConnectivity() error
}

// CypherDriver struct
type CypherDriver struct {
	db  *neoism.Database
	env string
}

//NewCypherDriver instantiate driver
func NewCypherDriver(db *neoism.Database, env string) CypherDriver {
	return CypherDriver{db, env}
}

// CheckConnectivity tests neo4j by running a simple cypher query
func (pcw CypherDriver) CheckConnectivity() error {
	results := []struct {
		ID int
	}{}
	query := &neoism.CypherQuery{
		Statement: "MATCH (x) RETURN ID(x) LIMIT 1",
		Result:    &results,
	}
	err := pcw.db.Cypher(query)
	log.Debugf("CheckConnectivity results:%+v  err: %+v", results, err)
	return err
}

type neoReadStruct struct {
	Uuid       string     `json:"uuid"`
	Types      []string   `json:"types"`
	Identifier Identifier `json:"identifier"`
}

type neoResultStrunct struct {
	Rs []neoReadStruct
}

func (pcw CypherDriver) ReadByConceptId(identifiers []string) (concordances Concordances, found bool, err error) {
	concordances = Concordances{}
	results := []neoResultStrunct{}
	query := &neoism.CypherQuery{
		Statement: `
		MATCH (p:Concept)<-[:IDENTIFIES]-(i:Identifier)
		WHERE p.uuid in {identifiers}
		RETURN collect({uuid:p.uuid, types:labels(p), Identifier:{authority:i.authority, identifierValue:i.value}}) as rs
		`,
		Parameters: neoism.Props{"identifiers": identifiers},
		Result:     &results,
	}

	return processCypherQueryToConcordances(pcw, query, &results)
}

func (pcw CypherDriver) ReadByAuthority(authority string, identifierValues []string) (concordances Concordances, found bool, err error) {
	concordances = Concordances{}
	results := []neoResultStrunct{}
	query := &neoism.CypherQuery{
		Statement: `
		MATCH (p:Concept)<-[:IDENTIFIES]-(i:Identifier)
		WHERE i.value in {identifierValues} AND i.authority = {authority}
		RETURN collect({uuid:p.uuid, types:labels(p), Identifier:{authority:i.authority, identifierValue:i.value}}) as rs
		`,
		Parameters: neoism.Props{
			"identifierValues": identifierValues,
			"authority":        authority,
		},
		Result: &results,
	}
	return processCypherQueryToConcordances(pcw, query, &results)
}

func processCypherQueryToConcordances(pcw CypherDriver, q *neoism.CypherQuery, results *[]neoResultStrunct) (concordances Concordances, found bool, err error) {
	err = pcw.db.Cypher(q)
	if err != nil {
		log.Errorf("Error looking up Concordances with query %s from neoism: %+v\n", q.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore for identifier:")
	}

	if (len(*results)) == 0 {
		return Concordances{}, false, nil
	}

	concordances = neoReadStructToConcordances(&(*results)[0].Rs, pcw.env)
	log.Debugf("Returning %v", concordances)
	return concordances, true, nil
}

func neoReadStructToConcordances(neo *[]neoReadStruct, env string) (concordances Concordances) {
	concordances = Concordances{
		Concordance: make([]Concordance, len(*neo)),
	}
	for i, neoCon := range *neo {
		var con = Concordance{}
		var concept = Concept{}
		concept.ID = neoCon.Uuid
		concept.APIURL = mapper.APIURL(neoCon.Uuid, neoCon.Types, env)
		con.Concept = concept
		con.Identifier = neoCon.Identifier
		concordances.Concordance[i] = con
	}
	return concordances
}
