package concordances

import (
	"fmt"

	"github.com/Financial-Times/neo-model-utils-go/mapper"
	log "github.com/Sirupsen/logrus"
	"github.com/jmcvetta/neoism"
)

// Driver interface
type Driver interface {
	ReadByConceptID(id string) (concordances Concordances, found bool, err error)
	ReadByAuthority(authority string, id string) (concordances Concordances, found bool, err error)
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
	UUID          string        `json:"uuid"`
	Types         []string      `json:"types"`
	NeoIdentifier neoIdentifier `json:"neoIdentifier"`
}

type neoIdentifier struct {
	Labels []string `json:"labels"`
	Value  string   `json:"value"`
}

type neoResultStrunct struct {
	Rs []neoReadStruct
}

func (pcw CypherDriver) ReadByConceptID(identifier string) (concordances Concordances, found bool, err error) {
	concordances = Concordances{}
	results := []neoResultStrunct{}
	query := &neoism.CypherQuery{
		Statement: `
		MATCH (p:Concept)<-[:IDENTIFIES]-(i:UPPIdentifier{value:{identifier}})
		MATCH (p:Concept)<-[:IDENTIFIES]-(ids:Identifier)
		RETURN collect({uuid:p.uuid, types:labels(p), neoIdentifier:{labels:labels(ids), value:ids.value}}) as rs
		`,
		Parameters: neoism.Props{"identifier": identifier},
		Result:     &results,
	}

	return processCypherQueryToConcordances(pcw, query, &results)
}

func (pcw CypherDriver) ReadByAuthority(authority string, identifierValue string) (concordances Concordances, found bool, err error) {
	concordances = Concordances{}
	results := []neoResultStrunct{}

	identifierLabel := mapAuthorityToIdentifierLabel(authority)

	if identifierLabel == "" {
		return Concordances{}, false, nil
	}

	readByAuthorityQueryStatement := fmt.Sprintf(`
		MATCH (p:Concept)<-[:IDENTIFIES]-(i:%s{value:{identifierValue}})
		MATCH (p:Concept)<-[:IDENTIFIES]-(ids:Identifier)
		RETURN collect({uuid:p.uuid, types:labels(p), neoIdentifier:{labels:labels(ids), value:ids.value}}) as rs
		`, identifierLabel)

	query := &neoism.CypherQuery{
		Statement: readByAuthorityQueryStatement,
		Parameters: neoism.Props{
			"identifierValue": identifierValue,
			"authority":       authority,
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

	if (len(concordances.Concordance)) == 0 {
		return Concordances{}, false, nil
	}

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
		concept.ID = mapper.IDURL(neoCon.UUID)
		concept.APIURL = mapper.APIURL(neoCon.UUID, neoCon.Types, env)
		con.Concept = concept
		con.Identifier = Identifier{Authority: mapNeoLabelsToAuthorityValue(neoCon.NeoIdentifier.Labels), IdentifierValue: neoCon.NeoIdentifier.Value}
		concordances.Concordance[i] = con
	}
	return concordances
}

func mapNeoLabelsToAuthorityValue(labelNames []string) (authority string) {
	for _, label := range labelNames {
		switch label {
		case TME_ID_NODE_LABEL:
			return TME_AUTHORITY
		case FS_ID_NODE_LABEL:
			return FS_AUTHORITY
		case UP_ID_NODE_LABEL:
			return UP_AUTHORITY
		case LEI_ID_NODE_LABEL:
			return LEI_AUTHORITY
		}
	}
	return ""
}

func mapAuthorityToIdentifierLabel(authority string) (label string) {
	switch authority {
	case TME_AUTHORITY:
		return TME_ID_NODE_LABEL
	case FS_AUTHORITY:
		return FS_ID_NODE_LABEL
	case UP_AUTHORITY:
		return UP_ID_NODE_LABEL
	case LEI_AUTHORITY:
		return LEI_ID_NODE_LABEL
	}
	return ""
}

const TME_AUTHORITY = "http://api.ft.com/system/FT-TME"
const FS_AUTHORITY = "http://api.ft.com/system/FACTSET"
const UP_AUTHORITY = "http://api.ft.com/system/UPP"
const LEI_AUTHORITY = "http://api.ft.com/system/LEI"

const TME_ID_NODE_LABEL = "TMEIdentifier"
const FS_ID_NODE_LABEL = "FactsetIdentifier"
const UP_ID_NODE_LABEL = "UPPIdentifier"
const LEI_ID_NODE_LABEL = "LegalEntityIdentifier"
