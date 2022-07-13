package zql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type DebugHandler struct {
}

func (d *DebugHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var q *Query
	var err error
	switch r.Header.Get("Content-Type") {
	case "application/json":
		q, err = d.parseJsonQuery(r.Body)
	case "plain/text":
		q, err = d.runExprQuery(r.Body)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := d.handleQuery(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)
}

func (d *DebugHandler) handleQuery(q *Query) ([]byte, error) {
	sqlt := NewSQLThesaurus("postgres", nil)
	cond, args, err := sqlt.ToSQL(q, true, true)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(nil)
	b.WriteString(fmt.Sprintf("WherePart: %s\n", cond))
	b.WriteString(fmt.Sprintf("Args: %v\n", args))
	return b.Bytes(), nil
}

func (d *DebugHandler) runExprQuery(r io.Reader) (*Query, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	q, err := Run(string(body), NewSyntaxV1())
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (d *DebugHandler) parseJsonQuery(r io.Reader) (*Query, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	q := new(Query)
	err = json.Unmarshal(body, q)
	if err != nil {
		return nil, err
	}
	return q, nil
}
