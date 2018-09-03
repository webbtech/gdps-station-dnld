package validate

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-fs-dwnld/model"
	"github.com/stretchr/testify/suite"
)

const (
	date       = "2018-05-01"
	stationID  = "d03224a7-f1df-4863-bcaa-5c6e61af11fc"
	timeFormat = "2006-01-02"
)

// UnitSuite struct
type UnitSuite struct {
	suite.Suite
	request *model.Request
}

// SetupTest method
func (suite *UnitSuite) SetupTest() {
	os.Setenv("Stage", "test")
	dte, err := time.Parse(timeFormat, date)
	suite.request = &model.Request{
		Date:      dte,
		StationID: stationID,
	}
	suite.NoError(err)
	suite.IsType(new(model.Request), suite.request)
}

// TestDate method
func (suite *UnitSuite) TestDate() {
	dte, err := Date(date)
	suite.NoError(err)
	suite.IsType(time.Time{}, dte)
}

// TestRequestInput method
func (suite *UnitSuite) TestRequestInput() {
	req := &model.RequestInput{
		Date:      date,
		StationID: stationID,
	}
	res, err := RequestInput(req)
	suite.NoError(err)
	suite.IsType(&model.Request{}, res)
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
