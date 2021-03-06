package goforce

import "testing"

import "fmt"

func TestInsert(t *testing.T) {
	force := NewClientTest()
	in := &InsertData{
		SObject: "Account",
		Fields: map[string]interface{}{
			"Name": "Testing",
		},
	}
	resp, err := force.Insert(in)
	if err != nil {
		t.Errorf("Error : %v", err)
	}
	fmt.Println(resp)
}

func TestUpdate(t *testing.T) {
	force := NewClientTest()
	in := &UpdateData{
		SObject: "Account",
		ID:      "0010w000009kJf8AAE",
		Fields: map[string]interface{}{
			"Name": "Updaters",
		},
	}
	err := force.Update(in)
	if err != nil {
		t.Errorf("Error : %v", err)
	}
}

func TestDelete(t *testing.T) {
	force := NewClientTest()
	del := &DeleteData{
		SObject: "Account",
		ID:      "0010w000009kJf8AAE",
	}
	err := force.Delete(del)
	if err != nil {
		t.Errorf("Error : %v", err)
	}
}
