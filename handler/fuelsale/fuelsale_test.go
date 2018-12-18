package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	// main "github.com/pulpfree/gdps-fs-dwnld/handler/fuelsale"
	"github.com/stretchr/testify/assert"
)

/*func TestHeartbeat(t *testing.T) {
	// suppose the use of convenience method of http.Get would work here as well...
	// req, err := http.NewRequest("GET", "http://127.0.0.1:3000/fuelsale", strings.NewReader(""))
	req, err := http.NewRequest("GET", "http://127.0.0.1:3000/fuelsale", nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("req: %+v\n", req)
	w := httptest.NewRecorder()
	fmt.Printf("w: %+v\n", w)
	fmt.Printf("w: %+v\n", w.Code)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("resp: %+v\n", resp.Body)
	fmt.Printf("body: %+v\n", body)

}*/

// const defaultsFilePath = "../config/defaults.yaml"

// UnitSuite struct
/*type UnitSuite struct {
	suite.Suite
	// client  *Client
	// cfg     *config.Config
	// request *model.Request
}

// SetupTest method
func (suite *UnitSuite) SetupTest() {
	// os.Setenv("Stage", "test")
	// os.Setenv("ConfigDefaults", "../config/defaults.yaml")
	fmt.Println("In setupTest")
}

// TestFuelSales method
func (suite *UnitSuite) TestFuelSales() {
	e := os.Getenv("ConfigDefaults")
	fmt.Printf("e in Test: %s\n", e)
	// res, err := suite.client.FuelSales()
	// suite.NoError(err)
	// suite.IsType(new(model.FuelSales), res)
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}*/

/*func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	fmt.Printf("TestMain called: %+v\n", m)
	const defaultsFilePath = "../config/defaults.yaml"
	os.Exit(m.Run())
}*/

func TestHandler(t *testing.T) {

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: "Paul"},
			expect:  "Hello Paul",
			err:     nil,
		},
	}

	for _, test := range tests {
		response, err := HandleRequest(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
