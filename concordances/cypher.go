package concordances

import (
	"fmt"

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

type neoReadIdentifier struct {
	authority       string
	identifierValue string
}

type neoReadStruct struct {
	uuid  string
	types []string
}

func (pcw CypherDriver) ReadByConceptId(identifiers []string) (concordances Concordances, found bool, err error) {
	concordances = Concordances{}
	results := []struct {
		Rs []neoReadStruct
	}{}
	query := &neoism.CypherQuery{
		Statement: `MATCH (p:Concept) where p.uuid in ["0b79b770-e426-334a-b231-f8f37d9a6678", "8138ca3f-b80d-3ef8-ad59-6a9b6ea5f15e"]
					OPTIONAL MATCH (p)<-[:IDENTIFIES]-(i:Identifier)
					RETURN p.uuid as uuid, labels(p) as types, collect(i) as identifiers  `,
		Parameters: neoism.Props{"identifier": identifiers},
		Result:     &results,
	}

	err = pcw.db.Cypher(query)
	fmt.Printf("neo:", err)
	if err != nil {
		log.Errorf("Error looking up uuid %s with query %s from neoism: %+v\n", identifiers, query.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore for identifier: %s", identifiers)
	}

	log.Debugf("CypherResult ReadConcordance for identifier: %s was: %+v", identifiers, results)

	fmt.Printf("RESULTS:", results[0])
	fmt.Printf("RESULTS1:", len(results))
	fmt.Printf("RESULTS2:", len(results[0].Rs))
	if (len(results)) == 0 || len(results[0].Rs) == 0 {
		fmt.Printf("ARGH")
		return Concordances{}, false, nil
	}

	concordances = neoReadStructToConcordances(results[0].Rs)
	log.Debugf("Returning %v", concordances)
	return concordances, true, nil
}

func neoReadStructToConcordances(neo []neoReadStruct) Concordances {
	fmt.Printf("IN FUNC:", neo)
	//TODO
	return Concordances{}
}
