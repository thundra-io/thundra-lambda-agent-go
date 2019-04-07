package thundrardb

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/lib/pq"

	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/tracer"
)

type postgresqlIntegration struct{}

func (i *postgresqlIntegration) getOperationName(query string) string {
	querySplit := strings.Split(query, " ")
	operation := ""
	if len(querySplit) > 0 {
		operation = querySplit[0]
	}
	return operation
}

func (i *postgresqlIntegration) beforeCall(query string, span *tracer.RawSpan, dsn string) {
	span.ClassName = constants.ClassNames["POSTGRESQL"]
	span.DomainName = constants.DomainNames["DB"]

	operation := i.getOperationName(query)

	dbName := ""
	host := ""
	port := ""
	opts, err := parseOptsFromDsn(dsn)
	if err == nil {
		dbName = opts["dbname"]
		host = opts["host"]
		port = opts["port"]
	}

	// Set span tags
	tags := map[string]interface{}{
		constants.SpanTags["OPERATION_TYPE"]:          operationToType[strings.ToLower(operation)],
		constants.SpanTags["TRIGGER_OPERATION_NAMES"]: []string{application.FunctionName},
		constants.SpanTags["TRIGGER_DOMAIN_NAME"]:     constants.AwsLambdaApplicationDomain,
		constants.SpanTags["TRIGGER_CLASS_NAME"]:      constants.AwsLambdaApplicationClass,
		constants.SpanTags["TOPOLOGY_VERTEX"]:         true,
		constants.DBTags["DB_STATEMENT_TYPE"]:         strings.ToUpper(operation),
		constants.DBTags["DB_TYPE"]:                   "postgresql",
		constants.DBTags["DB_STATEMENT_TYPE"]:         strings.ToUpper(operation),
		constants.DBTags["DB_INSTANCE"]:               dbName,
		constants.DBTags["DB_HOST"]:                   host,
		constants.DBTags["DB_PORT"]:                   port,
	}

	if !config.MaskRDBStatement {
		tags[constants.DBTags["DB_STATEMENT"]] = query
	}

	span.Tags = tags
}

func (i *postgresqlIntegration) afterCall(query string, span *tracer.RawSpan, dsn string) {
	return
}

type values map[string]string

// parseOptsFromDsn parses options from dsn and initializes and returns dbParams
func parseOptsFromDsn(dsn string) (values, error) {
	var err error
	o := make(values)
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		dsn, err = pq.ParseURL(dsn)
		if err != nil {
			return nil, err
		}
	}

	if err := parseOpts(dsn, o); err != nil {
		return nil, err
	}

	return o, nil
}

// scanner implements a tokenizer for libpq-style option strings.
type scanner struct {
	s []rune
	i int
}

// newScanner returns a new scanner initialized with the option string s.
func newScanner(s string) *scanner {
	return &scanner{[]rune(s), 0}
}

// Next returns the next rune.
// It returns 0, false if the end of the text has been reached.
func (s *scanner) next() (rune, bool) {
	if s.i >= len(s.s) {
		return 0, false
	}
	r := s.s[s.i]
	s.i++
	return r, true
}

// SkipSpaces returns the next non-whitespace rune.
// It returns 0, false if the end of the text has been reached.
func (s *scanner) skipSpaces() (rune, bool) {
	r, ok := s.next()
	for unicode.IsSpace(r) && ok {
		r, ok = s.next()
	}
	return r, ok
}

// parseOpts parses the options from name and adds them to the values.
//
// The parsing code is based on conninfo_parse from libpq's fe-connect.c
func parseOpts(name string, o values) error {
	s := newScanner(name)

	for {
		var (
			keyRunes, valRunes []rune
			r                  rune
			ok                 bool
		)

		if r, ok = s.skipSpaces(); !ok {
			break
		}

		// Scan the key
		for !unicode.IsSpace(r) && r != '=' {
			keyRunes = append(keyRunes, r)
			if r, ok = s.next(); !ok {
				break
			}
		}

		// Skip any whitespace if we're not at the = yet
		if r != '=' {
			r, ok = s.skipSpaces()
		}

		// The current character should be =
		if r != '=' || !ok {
			return fmt.Errorf(`missing "=" after %q in connection info string"`, string(keyRunes))
		}

		// Skip any whitespace after the =
		if r, ok = s.skipSpaces(); !ok {
			// If we reach the end here, the last value is just an empty string as per libpq.
			o[string(keyRunes)] = ""
			break
		}

		if r != '\'' {
			for !unicode.IsSpace(r) {
				if r == '\\' {
					if r, ok = s.next(); !ok {
						return fmt.Errorf(`missing character after backslash`)
					}
				}
				valRunes = append(valRunes, r)

				if r, ok = s.next(); !ok {
					break
				}
			}
		} else {
		quote:
			for {
				if r, ok = s.next(); !ok {
					return fmt.Errorf(`unterminated quoted string literal in connection string`)
				}
				switch r {
				case '\'':
					break quote
				case '\\':
					r, _ = s.next()
					fallthrough
				default:
					valRunes = append(valRunes, r)
				}
			}
		}

		o[string(keyRunes)] = string(valRunes)
	}

	return nil
}
