package goforce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	objectEndpoint = "/sobjects/"
)

// InsertValue is the value that is returned when a
// record is inserted into Salesforce.
type InsertValue struct {
	Success bool            `json:"success"`
	ID      string          `json:"id"`
	Errors  []ResponseError `json:"errors"`
}

// Inserter provides the parameters needed insert a record.
//
// SObject is the Salesforce table name.  An example would be Account or Custom__c.
//
// Fields are the fields of the record that are to be inserted.  It is the
// callers responsbility to provide value fields and values.
type Inserter interface {
	SObject() string
	Fields() map[string]interface{}
}

// Updater provides the parameters needed to update a record.
//
// SObject is the Salesforce table name.  An example would be Account or Custom__c.
//
// ID is the Salesforce ID that will be updated.
//
// Fields are the fields of the record that are to be inserted.  It is the
// callers responsbility to provide value fields and values.
type Updater interface {
	Inserter
	ID() string
}

// Deleter provides the parameters needed to delete a record.
//
// SObject is the Salesforce table name.  An example would be Account or Custom__c.
//
// ID is the Salesforce ID to be deleted.
type Deleter interface {
	SObject() string
	ID() string
}

// Represents insert data that is needed to do insert operation
type InsertData struct {
	sobject string
	fields  map[string]interface{}
}

func (insert *InsertData) SObject() string {
	return insert.sobject
}

func (insert *InsertData) Fields() map[string]interface{} {
	return insert.fields
}

// Represents update data that is needed to do update operation
type UpdateData struct {
	sobject string
	id      string
	fields  map[string]interface{}
}

func (update *UpdateData) SObject() string {
	return update.sobject
}

func (update *UpdateData) ID() string {
	return update.id
}

func (update *UpdateData) Fields() map[string]interface{} {
	return update.fields
}

// Represents delete data that is needed to do delete operation
type DeleteData struct {
	sobject string
	id      string
}

func (delete *DeleteData) SObject() string {
	return delete.sobject
}

func (delete *DeleteData) ID() string {
	return delete.id
}

func (force *ForceAPI) Insert(inserter Inserter) (InsertValue, error) {
	request, err := force.insertRequest(inserter)

	if err != nil {
		return InsertValue{}, err
	}

	value, err := force.insertResponse(request)

	if err != nil {
		return InsertValue{}, err
	}

	return value, nil
}

func (force *ForceAPI) insertRequest(inserter Inserter) (*http.Request, error) {

	url := force.getServiceURL() + objectEndpoint + inserter.SObject()

	body, err := json.Marshal(inserter.Fields())
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	return request, nil
}

func (force *ForceAPI) insertResponse(request *http.Request) (InsertValue, error) {
	response, err := force.client.Do(request)

	if err != nil {
		return InsertValue{}, err
	}

	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		var insertErrs []ResponseError
		err = decoder.Decode(&insertErrs)
		var errMsg error
		if err == nil {
			for _, insertErr := range insertErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", insertErr.ErrorCode, insertErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		}

		return InsertValue{}, errMsg
	}

	var value InsertValue
	err = decoder.Decode(&value)
	if err != nil {
		return InsertValue{}, err
	}

	return value, nil
}

func (force *ForceAPI) Update(updater Updater) error {
	request, err := force.updateRequest(updater)

	if err != nil {
		return err
	}

	err = force.updateResponse(request)

	if err != nil {
		return err
	}

	return nil
}

func (force *ForceAPI) updateRequest(updater Updater) (*http.Request, error) {

	url := force.getServiceURL() + objectEndpoint + updater.SObject() + "/" + updater.ID()

	body, err := json.Marshal(updater.Fields())
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	return request, nil
}

func (force *ForceAPI) updateResponse(request *http.Request) error {
	response, err := force.client.Do(request)
	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		var insertErrs []ResponseError
		err = decoder.Decode(&insertErrs)
		var errMsg error
		if err == nil {
			for _, insertErr := range insertErrs {
				errMsg = fmt.Errorf("insert response err: %s: %s", insertErr.ErrorCode, insertErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("insert response err: %d %s", response.StatusCode, response.Status)
		}

		return errMsg
	}
	return nil
}

func (force *ForceAPI) Delete(deleter Deleter) error {
	request, err := force.deleteRequest(deleter)

	if err != nil {
		return err
	}

	err = force.deleteResponse(request)

	if err != nil {
		return err
	}

	return nil
}

func (force *ForceAPI) deleteRequest(deleter Deleter) (*http.Request, error) {

	url := force.getServiceURL() + objectEndpoint + deleter.SObject() + "/" + deleter.ID()

	request, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	return request, nil
}

func (force *ForceAPI) deleteResponse(request *http.Request) error {
	response, err := force.client.Do(request)
	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		var insertErrs []ResponseError
		err = decoder.Decode(&insertErrs)
		var errMsg error
		if err == nil {
			for _, insertErr := range insertErrs {
				errMsg = fmt.Errorf("delete response err: %s: %s", insertErr.ErrorCode, insertErr.Message)
			}
		} else {
			errMsg = fmt.Errorf("delete response err: %d %s", response.StatusCode, response.Status)
		}

		return errMsg
	}
	return nil
}
