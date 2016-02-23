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

func (pcw CypherDriver) ReadByConceptId(identifiers []string) (concordances Concordances, found bool, err error) {
	concordances = Concordances{}
	results := []struct{ Rs []neoReadStruct }{}
	query := &neoism.CypherQuery{
		Statement: `
		MATCH (p:Concept)<-[:IDENTIFIES]-(i:Identifier)
		WHERE p.uuid in {identifiers}
		RETURN collect({uuid:p.uuid, types:labels(p), Identifier:{authority:i.authority, identifierValue:i.value}}) as rs
		`,
		Parameters: neoism.Props{"identifiers": identifiers},
		Result:     &results,
	}

	err = pcw.db.Cypher(query)
	fmt.Printf("neo_err:%s\n", err)
	if err != nil {
		log.Errorf("Error looking up uuid %s with query %s from neoism: %+v\n", identifiers, query.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore for identifier: %s", identifiers)
	}

	log.Debugf("CypherResult ReadConcordance for identifier: %s was: %+v", identifiers, results)

	fmt.Printf("RESULTS:%s\n", results)
	if (len(results)) == 0 {
		fmt.Printf("ARGH\n")
		return Concordances{}, false, nil
	}

	concordances = neoReadStructToConcordances(results[0].Rs, pcw.env)
	log.Debugf("Returning %v", concordances)
	return concordances, true, nil
}

func (pcw CypherDriver) ReadByAuthority(authority string, identifierValues []string) (concordances Concordances, found bool, err error) {
	concordances = Concordances{}
	results := []struct{ Rs []neoReadStruct }{}
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

	err = pcw.db.Cypher(query)
	fmt.Printf("neo_err:%s\n", err)
	if err != nil {
		log.Errorf("Error looking up uuid %s with query %s from neoism: %+v\n", identifierValues, query.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore for identifierValues: %s", identifierValues)
	}

	log.Debugf("CypherResult ReadConcordance for identifierValues: %s was: %+v", identifierValues, results)

	fmt.Printf("RESULTS:%s\n", results)
	if (len(results)) == 0 {
		fmt.Printf("ARGH\n")
		return Concordances{}, false, nil
	}

	concordances = neoReadStructToConcordances(results[0].Rs, pcw.env)
	log.Debugf("Returning %v", concordances)
	return concordances, true, nil
}

func neoReadStructToConcordances(neo []neoReadStruct, env string) Concordances {
	var concordances = Concordances{
		Concordance: make([]Concordance, len(neo)),
	}
	for i, neoCon := range neo {
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
