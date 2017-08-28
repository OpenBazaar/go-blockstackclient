package blockstackclient

import (
	"testing"
	"bytes"
	"net/http"
	"path"
	"gx/ipfs/QmTEmsyNnckEq8rEfALfdhLHjrEHGoSGFDrAYReuetn7MC/go-net/context"
)

var testClient *BlockstackClient

type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb *ClosingBuffer) Close() (err error) {
	return
}

type mockHttpClient struct{}

func (m *mockHttpClient) Get(url string) (*http.Response, error){
	data := `{
	  "testuser": {
	    "profile": {
	      "account": [
		{
		  "identifier": "QmdHkAQeKJobghWES9exVUaqXCeMw8katQitnXDKWuKi1F",
		  "proofType": "http",
		  "@type": "Account",
		  "service": "openbazaar"
		}
	      ]
	    }
	  }
	}`
	_, username := path.Split(url)

	cb := &ClosingBuffer{bytes.NewBufferString(data)}
	resp := &http.Response{
		Body: cb,
	}
	if username != "testuser" {
		resp.StatusCode = http.StatusNotFound
	}
	return resp, nil
}

func init(){
	testClient = &BlockstackClient{
		resolverURL: "http://xyz.com/",
		httpClient: &mockHttpClient{},
	}
}

func TestBlockstackClient_Resolve(t *testing.T) {
	guid, err := testClient.Resolve(context.Background(), "testuser.id")
	if err != nil {
		t.Error(err)
	}
	if guid.Pretty() != "QmdHkAQeKJobghWES9exVUaqXCeMw8katQitnXDKWuKi1F" {
		t.Error("Returned invalid guid")
	}

	_, err = testClient.Resolve(context.Background(), "nonexistantuser")
	if err == nil {
		t.Error(err)
	}

}