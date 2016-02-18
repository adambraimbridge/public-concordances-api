package concordances

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/jmcvetta/neoism"
)

// Driver interface
type Driver interface {
	Read(id string) (concordance Concordance, found bool, err error)
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

type neoChangeEvent struct {
	StartedAt string
	EndedAt   string
}

type neoReadStruct struct {
	O struct {
		ID        string
		Types     []string
		LEICode   string
		PrefLabel string
		Labels    []string
	}
	Parent struct {
		ID        string
		Types     []string
		PrefLabel string
	}
	Ind struct {
		ID        string
		Types     []string
		PrefLabel string
	}
	Sub []struct {
		ID        string
		Types     []string
		PrefLabel string
	}
	PM []struct {
		M struct {
			ID           string
			Types        []string
			PrefLabel    string
			Title        string
			ChangeEvents []neoChangeEvent
		}
		P struct {
			ID        string
			Types     []string
			PrefLabel string
			Labels    []string
		}
	}
}

func (pcw CypherDriver) Read(uuid string) (concordance Concordance, found bool, err error) {
	concordance = Concordance{}
	results := []struct {
		Rs []neoReadStruct
	}{}
	query := &neoism.CypherQuery{
		Statement:  `TODO`,
		Parameters: neoism.Props{"uuid": uuid},
		Result:     &results,
	}
	err = pcw.db.Cypher(query)
	if err != nil {
		log.Errorf("Error looking up uuid %s with query %s from neoism: %+v\n", uuid, query.Statement, err)
		return Concordance{}, false, fmt.Errorf("Error accessing Concordance datastore for uuid: %s", uuid)
	}
	log.Debugf("CypherResult ReadConcordance for uuid: %s was: %+v", uuid, results)
	if (len(results)) == 0 || len(results[0].Rs) == 0 {
		return Concordance{}, false, nil
	} else if len(results) != 1 && len(results[0].Rs) != 1 {
		errMsg := fmt.Sprintf("Multiple concordances found with the same uuid:%s !", uuid)
		log.Error(errMsg)
		return Concordance{}, true, errors.New(errMsg)
	}
	concordance = neoReadStructToConcordance(results[0].Rs[0])
	log.Debugf("Returning %v", concordance)
	return concordance, true, nil
}

func neoReadStructToConcordance(neo neoReadStruct) Concordance {
	//TODO
	return Concordance{}
}
