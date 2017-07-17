package concordances

import (
	"fmt"

	"github.com/Financial-Times/neo-model-utils-go/mapper"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	log "github.com/Sirupsen/logrus"
	"github.com/jmcvetta/neoism"
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
	UUID          string        `json:"UUID"`
	Types         []string      `json:"TYPES"`
	NeoIdentifier neoIdentifier `json:"IDENTIFIERS"`
	PrefUUID      string        `json:"prefUUID"`
}

type neoIdentifier struct {
	Labels []string `json:"labels"`
	Value  string   `json:"value"`
}

type neoResultStrunct struct {
	Rs []neoReadStruct
}

func (pcw CypherDriver) ReadByConceptID(identifiers []string) (concordances Concordances, found bool, err error) {
	c, f, err := pcw.readByConceptIDNewModel(identifiers)
	if !f {
		c, f, err = pcw.readByConceptIDOldModel(identifiers)
	}
	return c, f, err
}

func (pcw CypherDriver) readByConceptIDNewModel(identifiers []string) (concordances Concordances, found bool, err error) {
	results := []neoReadStruct{}
	query := &neoism.CypherQuery{
		Statement: `
		MATCH (p:Concept)-[:EQUIVALENT_TO]-(cn:Concept)
		WHERE p.uuid in {identifiers}
		MATCH (cn)-[:EQUIVALENT_TO]-(cnn:Concept)
		RETURN cn.prefUUID as prefUUID, cnn.uuid AS UUID, labels(cnn) AS TYPES, {labels:collect(cnn.authority), value:cnn.authorityValue} as IDENTIFIERS
		`,
		Parameters: neoism.Props{"identifiers": identifiers},
		Result:     &results,
	}

	err = pcw.conn.CypherBatch([]*neoism.CypherQuery{query})
	if err != nil {
		log.Errorf("Error looking up Concordances with query %s from neoism: %+v\n", query.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore for identifier:")
	}

	log.Info(results)

	if (len(results)) == 0 {
		return Concordances{}, false, nil
	}

	concordances = Concordances{
		Concordance: []Concordance{},
	}

	for _, neoCon := range results {
		log.Debug(neoCon)
		var con = Concordance{}
		var concept = Concept{}
		concept.ID = mapper.IDURL(neoCon.PrefUUID)
		concept.APIURL = mapper.APIURL(neoCon.PrefUUID, neoCon.Types, pcw.env)
		con.Concept = concept
		con.Identifier = Identifier{Authority: mapNeoLabelsToAuthorityValue(neoCon.NeoIdentifier.Labels), IdentifierValue: neoCon.NeoIdentifier.Value}
		concordances.Concordance = append(concordances.Concordance, con)
	}

	if (len(concordances.Concordance)) == 0 {
		return Concordances{}, false, nil
	}

	log.Debugf("Returning %v", concordances)
	return concordances, true, nil
}

func (pcw CypherDriver) readByConceptIDOldModel(identifiers []string) (concordances Concordances, found bool, err error) {
	results := []neoReadStruct{}
	query := &neoism.CypherQuery{
		Statement: `
		MATCH (p:Concept)<-[:IDENTIFIES]-(i:UPPIdentifier)
		WHERE i.value in {identifiers}
		MATCH (p:Concept)<-[:IDENTIFIES]-(ids:Identifier)
		WHERE NOT ids:UPPIdentifier
		RETURN p.uuid as prefUUID, p.uuid AS UUID, labels(p) AS TYPES, {labels:labels(ids), value:ids.value} as IDENTIFIERS
		`,
		Parameters: neoism.Props{"identifiers": identifiers},
		Result:     &results,
	}
	return processCypherQueryToConcordances(pcw, query, &results)
}

func (pcw CypherDriver) ReadByAuthority(authority string, identifierValues []string) (concordances Concordances, found bool, err error) {
	c, f, err := pcw.readByAuthorityNewModel(authority, identifierValues)
	if !f {
		c, f, err = pcw.readByAuthorityOldModel(authority, identifierValues)
	}
	return c, f, err
}

func (pcw CypherDriver) readByAuthorityNewModel(authority string, identifierValues []string) (concordances Concordances, found bool, err error) {
	log.Debug("readByAuthorityNewModel")
	concordances = Concordances{}
	results := []neoReadStruct{}

	authorityProperty := mapAuthorityToAuthorityProperty(authority)
	log.Info(authorityProperty)
	if authorityProperty == "" {
		return Concordances{}, false, nil
	}

	readByAuthorityQueryStatement := `
		MATCH (p:Concept)-[:EQUIVALENT_TO]-(cn:Concept)
		WHERE p.authorityValue in {identifierValues} AND p.authority = {authority}
		RETURN cn.prefUUID as prefUUID, cn.prefUUID AS UUID, labels(cn) AS TYPES, {labels:collect(p.authority), value:p.authorityValue} as IDENTIFIERS
		`

	query := &neoism.CypherQuery{
		Statement: readByAuthorityQueryStatement,
		Parameters: neoism.Props{
			"identifierValues": identifierValues,
			"authority":        authorityProperty,
		},
		Result: &results,
	}

	err = pcw.conn.CypherBatch([]*neoism.CypherQuery{query})
	if err != nil {
		log.Errorf("Error looking up Concordances with query %s from neoism: %+v\n", query.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore for identifier:")
	}

	log.Info(results)

	if (len(results)) == 0 {
		return Concordances{}, false, nil
	}

	concordances = Concordances{
		Concordance: []Concordance{},
	}

	for _, neoCon := range results {
		log.Info(neoCon)
		// Each record is now two identifiers, one UPP and one other.
		var con = Concordance{}
		var concept = Concept{}
		concept.ID = mapper.IDURL(neoCon.PrefUUID)
		concept.APIURL = mapper.APIURL(neoCon.PrefUUID, neoCon.Types, pcw.env)
		con.Concept = concept
		con.Identifier = Identifier{Authority: mapNeoLabelsToAuthorityValue(neoCon.NeoIdentifier.Labels), IdentifierValue: neoCon.NeoIdentifier.Value}
		concordances.Concordance = append(concordances.Concordance, con)
	}

	if (len(concordances.Concordance)) == 0 {
		return Concordances{}, false, nil
	}

	log.Debugf("Returning %v", concordances)
	return concordances, true, nil
}

func (pcw CypherDriver) readByAuthorityOldModel(authority string, identifierValues []string) (concordances Concordances, found bool, err error) {
	log.Debug("readByAuthorityOldModel")
	concordances = Concordances{}
	results := []neoReadStruct{}

	identifierLabel := mapAuthorityToIdentifierLabel(authority)

	if identifierLabel == "" {
		return Concordances{}, false, nil
	}

	readByAuthorityQueryStatement := fmt.Sprintf(`
		MATCH (p:Concept)<-[:IDENTIFIES]-(i:%s)
		WHERE i.value in {identifierValues}
		RETURN p.uuid as prefUUID, p.uuid AS UUID, labels(p) AS TYPES, {labels:labels(i), value:i.value} as IDENTIFIERS
		`, identifierLabel)

	query := &neoism.CypherQuery{
		Statement: readByAuthorityQueryStatement,
		Parameters: neoism.Props{
			"identifierValues": identifierValues,
			"authority":        authority,
		},
		Result: &results,
	}
	return processCypherQueryToConcordances(pcw, query, &results)
}

func processCypherQueryToConcordances(pcw CypherDriver, q *neoism.CypherQuery, results *[]neoReadStruct) (concordances Concordances, found bool, err error) {
	err = pcw.conn.CypherBatch([]*neoism.CypherQuery{q})
	if err != nil {
		log.Errorf("Error looking up Concordances with query %s from neoism: %+v\n", q.Statement, err)
		return Concordances{}, false, fmt.Errorf("Error accessing Concordance datastore for identifier:")
	}

	if (len(*results)) == 0 {
		return Concordances{}, false, nil
	}

	concordances = neoReadStructToConcordances(results, pcw.env)

	if (len(concordances.Concordance)) == 0 {
		return Concordances{}, false, nil
	}

	log.Debugf("Returning %v", concordances)
	return concordances, true, nil
}

func neoReadStructToConcordances(neo *[]neoReadStruct, env string) (concordances Concordances) {
	log.Debug("Running the old model")
	concordances = Concordances{
		Concordance: []Concordance{},
	}
	for _, neoCon := range *neo {
		var con = Concordance{}
		var concept = Concept{}
		concept.ID = mapper.IDURL(neoCon.UUID)
		concept.APIURL = mapper.APIURL(neoCon.UUID, neoCon.Types, env)
		con.Concept = concept
		con.Identifier = Identifier{Authority: mapNeoLabelsToAuthorityValue(neoCon.NeoIdentifier.Labels), IdentifierValue: neoCon.NeoIdentifier.Value}
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
