package concordances

import (
	"fmt"

	"github.com/Financial-Times/neo-model-utils-go/mapper"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
	log "github.com/sirupsen/logrus"
)

// Driver interface
type Driver interface {
	ReadByConceptID(ids []string) (concordances Concordances, found bool, err error)
	ReadByAuthority(authority string, ids []string) (concordances Concordances, found bool, err error)
	CheckConnectivity() error
}

// CypherDriver struct
type CypherDriver struct {
	conn neoutils.NeoConnection
	env  string
}

//NewCypherDriver instantiate driver
func NewCypherDriver(conn neoutils.NeoConnection, env string) CypherDriver {
	return CypherDriver{conn, env}
}

// CheckConnectivity tests neo4j by running a simple cypher query
func (pcw CypherDriver) CheckConnectivity() error {
	return neoutils.Check(pcw.conn)
}

type neoReadStruct struct {
	UUID        string        `json:"UUID"`
	Types       []string      `json:"types"`
	Identifiers neoIdentifier `json:"identifiers"`
}

type neoIdentifier struct {
	Labels []string `json:"labels"`
	Value  string   `json:"value"`
}

func (pcw CypherDriver) ReadByConceptID(identifiers []string) (concordances Concordances, found bool, err error) {
	results := []neoReadStruct{}
	query := &neoism.CypherQuery{
		Statement: `
		MATCH (p:Concept)<-[:IDENTIFIES]-(i:UPPIdentifier)
		WHERE i.value in {identifiers}
		MATCH (p)-[:IDENTIFIES]-(ids:Identifier)
		OPTIONAL MATCH (p)-[:EQUIVALENT_TO]->(canonical:Concept)
		OPTIONAL MATCH (leafNode:Concept)-[:EQUIVALENT_TO]->(canonical)
		OPTIONAL MATCH (leafNode)<-[:IDENTIFIES]-(leafId:Identifier)
		WITH COALESCE(canonical.prefUUID, p.uuid) AS UUID, COALESCE(labels(canonical), labels(p)) AS types, COALESCE(leafId, ids) as nodeIds
		RETURN DISTINCT UUID, types, {labels:labels(nodeIds), value:nodeIds.value} as identifiers
        `,
		Parameters: neoism.Props{"identifiers": identifiers},
		Result:     &results,
	}

	err = pcw.conn.CypherBatch([]*neoism.CypherQuery{query})
	if err != nil {
		log.Errorf("Error looking up Concordances with query %s from neoism: %+v\n", query.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore for identifier: %v", identifiers)
	}

	if (len(results)) == 0 {
		return Concordances{}, false, nil
	}

	concordances = Concordances{
		Concordance: []Concordance{},
	}

	return processCypherQueryToConcordances(pcw, query, results)

}

func (pcw CypherDriver) ReadByAuthority(authority string, identifierValues []string) (concordances Concordances, found bool, err error) {
	results := []neoReadStruct{}

	authorityLabel := mapAuthorityToIdentifierLabel(authority)
	if authorityLabel == "" {
		return Concordances{}, false, nil
	}

	query := &neoism.CypherQuery{
		Statement: fmt.Sprintf(`
		MATCH (p:Concept)<-[:IDENTIFIES]-(ids:Identifier:%s)
 		WHERE ids.value in {identifierValues}
		OPTIONAL MATCH (p)-[:EQUIVALENT_TO]->(canonical:Concept)
		RETURN COALESCE(canonical.prefUUID, p.uuid) AS UUID, COALESCE(labels(canonical), labels(p)) AS types, {labels:labels(ids), value:ids.value} as identifiers`, authorityLabel),

		Parameters: neoism.Props{
			"identifierValues": identifierValues,
		},
		Result: &results,
	}

	err = pcw.conn.CypherBatch([]*neoism.CypherQuery{query})
	if err != nil {
		log.Errorf("Error looking up Concordances with query %s from neoism: %+v\n", query.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore for identifier: %v", identifierValues)
	}

	if (len(results)) == 0 {
		return Concordances{}, false, nil
	}

	concordances = Concordances{
		Concordance: []Concordance{},
	}

	return processCypherQueryToConcordances(pcw, query, results)
}

func processCypherQueryToConcordances(pcw CypherDriver, q *neoism.CypherQuery, results []neoReadStruct) (concordances Concordances, found bool, err error) {
	err = pcw.conn.CypherBatch([]*neoism.CypherQuery{q})
	if err != nil {
		log.Errorf("Error looking up Concordances with query %s from neoism: %+v\n", q.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore")
	}

	concordances = neoReadStructToConcordances(results, pcw.env)

	if (len(concordances.Concordance)) == 0 {
		return Concordances{}, false, nil
	}
	return concordances, true, nil
}

func neoReadStructToConcordances(neo []neoReadStruct, env string) (concordances Concordances) {
	concordances = Concordances{
		Concordance: []Concordance{},
	}
	for _, neoCon := range neo {
		var con = Concordance{}
		var concept = Concept{}

		log.Debug(neoCon)

		concept.ID = mapper.IDURL(neoCon.UUID)
		concept.APIURL = mapper.APIURL(neoCon.UUID, neoCon.Types, env)
		con.Identifier = Identifier{Authority: mapNeoLabelsToAuthorityValue(neoCon.Identifiers.Labels), IdentifierValue: neoCon.Identifiers.Value}

		con.Concept = concept
		concordances.Concordance = append(concordances.Concordance, con)
	}
	return concordances
}

func mapNeoLabelsToAuthorityValue(labelNames []string) (authority string) {
	for _, label := range labelNames {
		switch label {
		// Old style node label lookup
		case TME_ID_NODE_LABEL:
			return TME_AUTHORITY
		case FS_ID_NODE_LABEL:
			return FS_AUTHORITY
		case UP_ID_NODE_LABEL:
			return UP_AUTHORITY
		case LEI_ID_NODE_LABEL:
			return LEI_AUTHORITY
		case SL_ID_NODE_LABEL:
			return SL_AUTHORITY

		// New style authority properties
		case FS_AUTHORITY_PROPERTY:
			return FS_AUTHORITY
		case UP_AUTHORITY_PROPERTY:
			return UP_AUTHORITY
		case SL_AUTHORITY_PROPERTY:
			return SL_AUTHORITY
		case TME_AUTHORITY_PROPERTY:
			return TME_AUTHORITY
		}
	}
	return ""
}

func mapAuthorityToAuthorityProperty(authority string) string {
	switch authority {
	case TME_AUTHORITY:
		return TME_AUTHORITY_PROPERTY
	case FS_AUTHORITY:
		return FS_AUTHORITY_PROPERTY
	case SL_AUTHORITY:
		return SL_AUTHORITY_PROPERTY
	default:
		return ""
	}
}

func mapAuthorityToIdentifierLabel(authority string) (label string) {
	switch authority {
	case UP_AUTHORITY:
		return UP_ID_NODE_LABEL
	case SL_AUTHORITY:
		return SL_ID_NODE_LABEL
	case TME_AUTHORITY:
		return TME_ID_NODE_LABEL
	case FS_AUTHORITY:
		return FS_ID_NODE_LABEL
	case LEI_AUTHORITY:
		return LEI_ID_NODE_LABEL
	}
	return ""
}

const TME_AUTHORITY = "http://api.ft.com/system/FT-TME"
const FS_AUTHORITY = "http://api.ft.com/system/FACTSET"
const UP_AUTHORITY = "http://api.ft.com/system/UPP"
const LEI_AUTHORITY = "http://api.ft.com/system/LEI"
const SL_AUTHORITY = "http://api.ft.com/system/SMARTLOGIC"

const TME_ID_NODE_LABEL = "TMEIdentifier"
const FS_ID_NODE_LABEL = "FactsetIdentifier"
const UP_ID_NODE_LABEL = "UPPIdentifier"
const LEI_ID_NODE_LABEL = "LegalEntityIdentifier"
const SL_ID_NODE_LABEL = "SmartlogicIdentifier"

const TME_AUTHORITY_PROPERTY = "TME"
const FS_AUTHORITY_PROPERTY = "FACTSET"
const UP_AUTHORITY_PROPERTY = "UPP"
const SL_AUTHORITY_PROPERTY = "Smartlogic"
