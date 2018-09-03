package graphql

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-fs-dwnld/config"
	"github.com/pulpfree/gdps-fs-dwnld/model"
	"github.com/stretchr/testify/suite"
)

const (
	date             = "2018-05-01"
	defaultsFilePath = "../config/defaults.yaml"
	stationID        = "d03224a7-f1df-4863-bcaa-5c6e61af11fc"
	timeFormat       = "2006-01-02"
)

// UnitSuite struct
type UnitSuite struct {
	suite.Suite
	client  *Client
	cfg     *config.Config
	request *model.Request
	// report  *Report
}

// SetupTest method
func (suite *UnitSuite) SetupTest() {
	os.Setenv("Stage", "test")
	dte, err := time.Parse(timeFormat, date)
	req := &model.Request{
		Date:      dte,
		StationID: stationID,
	}
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFilePath}
	err = suite.cfg.Load()
	suite.NoError(err)
	suite.IsType(new(config.Config), suite.cfg)

	suite.client = New(req, suite.cfg)
	suite.NoError(err)
	suite.IsType(new(Client), suite.client)
}

// TestFuelSales method
func (suite *UnitSuite) TestFuelSales() {
	res, err := suite.client.FuelSales()
	suite.NoError(err)
	suite.IsType(new(model.FuelSales), res)
}

// TestFuelDelivery method
func (suite *UnitSuite) TestFuelDelivery() {
	res, err := suite.client.FuelDelivery()
	suite.NoError(err)
	suite.IsType(new(model.FuelDelivery), res)
}

// TestOverShortMonth method
func (suite *UnitSuite) TestOverShortMonth() {
	res, err := suite.client.OverShortMonth()
	suite.NoError(err)
	suite.IsType(new(model.OverShortMonth), res)
}

// TestOverShortAnnual method
func (suite *UnitSuite) TestOverShortAnnual() {
	res, err := suite.client.OverShortAnnual()
	suite.NoError(err)
	suite.IsType(new(model.OverShortAnnual), res)
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
