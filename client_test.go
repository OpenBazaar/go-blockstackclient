package blockstackclient

import (
	"testing"
	"bytes"
	"net/http"
	"path"
	"time"
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
		cache: make(map[string]CachedGuid),
		cacheLife: time.Minute,
	}
}

func TestBlockstackClient_Resolve(t *testing.T) {
	guid, err := testClient.Resolve("@testuser")
	if err != nil {
		t.Error(err)
	}
	if guid != "QmdHkAQeKJobghWES9exVUaqXCeMw8katQitnXDKWuKi1F" {
		t.Error("Returned invalid guid")
	}

	if _, ok := testClient.cache["testuser"]; !ok {
		t.Error("Client failed to cache response")
	}

	_, err = testClient.Resolve("@nonexistantuser")
	if err == nil {
		t.Error(err)
	}

	testClient.cache["cacheduser"] = CachedGuid{"abc", time.Now()}
	guid, err = testClient.Resolve("@cacheduser")
	if err != nil {
		t.Error(err)
	}
	if guid != "abc" {
		t.Error("Returned incorrect guid from cache")
	}

}

func TestFormatHandle(t *testing.T) {
	if formatHandle("@testuser") != "testuser" {
		t.Error("Failed to correctly format handle")
	}
	if formatHandle("@TestUser") != "testuser" {
		t.Error("Failed to correctly format handle")
	}
	if formatHandle("testuser") != "testuser" {
		t.Error("Failed to correctly format handle")
	}
}

func TestBlockstackClient_DeleteExpiredCache(t *testing.T) {
	testClient.cache["testexpired"] = CachedGuid{"abc", time.Now()}
	testClient.deleteExpiredCache()
	if _, ok := testClient.cache["testexpired"]; ok {
		t.Error("Client failed to delete expired cache")
	}
}
