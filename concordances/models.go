package concordances

// Concordances is a list of concordances wrapped like this for parity in the JSON currently produced
type Concordances struct {
	Concordance []Concordance `json:"concordances,omitempty"`
}

// Concept is a concept equivilant to a thing
/* This is equivilant to a Thing without  a prefLabel and probably at some point
should be converted to Thing but initially we arew only wanting parity with the existing
service
public class ConceptView {
    private String id;
    private String apiUrl;
}
*/
type Concept struct {
	ID     string `json:"id"`
	APIURL string `json:"apiUrl"` // self ?
}

// Concordance is the structure used for the people API
/* The following is currently defined in Java (e4b93668e32) but I think we should be removing profile
public class Concordance {
    private ConceptView concept;
    private Identifier identifier;
}
*/
type Concordance struct {
	Concept    Concept    `json:"concept,omitempty"`
	Identifier Identifier `json:"indentifier,omitempty"`
}

// Identifier identifies the concept with alternative identity
/*public class Identifier {

  private String authority;
  private String identifierValue;
*/
type Identifier struct {
	Authority       string `json:"authority"`
	IdentifierValue string `json:"identifierValue"`
}
