package graphql

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/machinebox/graphql"
	"github.com/pulpfree/gdps-fs-dwnld/config"
	"github.com/pulpfree/gdps-fs-dwnld/model"

	log "github.com/sirupsen/logrus"
)

// Client struct
type Client struct {
	hdrs    http.Header
	client  *graphql.Client
	request *model.Request
}

const timeLongFrmt = "2006-01-02"

// New graphql client
func New(req *model.Request, cfg *config.Config, authToken string) (c *Client) {

	hdrs := http.Header{}
	if len(authToken) > 0 {
		hdrs.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	c = &Client{
		client:  graphql.NewClient(cfg.GraphqlURI),
		hdrs:    hdrs,
		request: req,
	}

	return c
}

// FuelSales method
func (c *Client) FuelSales() (rpt *model.FuelSales, err error) {

	req := graphql.NewRequest(`
    query ($date: String!, $stationID: String!) {
      station(stationID: $stationID) {
        id
        name
      }
      fuelSaleMonth(date: $date, stationID: $stationID) {
        fuelTypes
        stationSales {
          date
          sales
        }
        salesSummary
        salesTotal
      }
    }
  `)

	req.Var("date", formattedDate(c.request.Date))
	req.Var("stationID", c.request.StationID)
	req.Header = c.hdrs

	ctx := context.Background()
	err = c.client.Run(ctx, req, &rpt)
	if err != nil {
		log.Errorf("error running graphql client: %s", err.Error())
		return nil, err
	}

	rpt.Date = c.request.Date
	rpt.Report.FuelTypes = sortFuelTypes(rpt.Report.FuelTypes)

	return rpt, err
}

// FuelDelivery method
func (c *Client) FuelDelivery() (rpt *model.FuelDelivery, err error) {

	req := graphql.NewRequest(`
    query FuelDeliveryReport($date: String!, $stationID: String!) {
      station(stationID: $stationID) {
        id
        name
      }
      fuelDeliveryReport(date: $date, stationID: $stationID) {
        fuelTypes
        deliveries {
          data
          date
        }
        deliverySummary
      }
    }
  `)

	req.Var("date", formattedDate(c.request.Date))
	req.Var("stationID", c.request.StationID)
	req.Header = c.hdrs

	ctx := context.Background()
	err = c.client.Run(ctx, req, &rpt)
	if err != nil {
		log.Errorf("error running graphql client: %s", err.Error())
		return nil, err
	}

	rpt.Report.FuelTypes = sortFuelTypes(rpt.Report.FuelTypes)
	rpt.Date = c.request.Date

	return rpt, err
}

// OverShortMonth method
func (c *Client) OverShortMonth() (rpt *model.OverShortMonth, err error) {

	req := graphql.NewRequest(`
    query DipOSMonthReport($date: String!, $stationID: String!) {
      station(stationID: $stationID) {
        id
        name
      }
      dipOSMonthReport(date: $date, stationID: $stationID) {
        stationID
        fuelTypes
        period
        overShort {
          date
          data
        }
        overShortSummary
      }
    }
  `)

	req.Var("date", formattedDate(c.request.Date))
	req.Var("stationID", c.request.StationID)
	req.Header = c.hdrs

	ctx := context.Background()
	err = c.client.Run(ctx, req, &rpt)
	if err != nil {
		log.Errorf("error running graphql client: %s", err.Error())
		return nil, err
	}

	rpt.Report.FuelTypes = sortFuelTypes(rpt.Report.FuelTypes)
	rpt.Date = c.request.Date

	return rpt, err
}

// OverShortAnnual method
func (c *Client) OverShortAnnual() (rpt *model.OverShortAnnual, err error) {

	req := graphql.NewRequest(`
    query DipOSAnnualReport($date: String!, $stationID: String!) {
      station(stationID: $stationID) {
        id
        name
      }
      dipOSAnnualReport(date: $date, stationID: $stationID) {
        fuelTypes
        year
        months
        summary
      }
    }
  `)

	req.Var("date", formattedDate(c.request.Date))
	req.Var("stationID", c.request.StationID)
	req.Header = c.hdrs

	ctx := context.Background()
	err = c.client.Run(ctx, req, &rpt)
	if err != nil {
		log.Errorf("error running graphql client: %s", err.Error())
		return nil, err
	}

	rpt.Report.FuelTypes = sortFuelTypes(rpt.Report.FuelTypes)
	rpt.Date = c.request.Date

	return rpt, err
}

// FuelSalesList method
func (c *Client) FuelSalesList() (rpt *model.FuelSalesList, err error) {

	req := graphql.NewRequest(`
    query FuelSaleListReport($date: String!) {
      fuelSaleListReport(date: $date) {
        periodHeader {
          yearWeek
          startDate
          endDate
          week
        }
        periodSales {
          fuelPrices
          periods {
            dates
            fuelSales {
              NL
              DSL
            }
          }
          stationID
          stationName
          stationTotal {
            NL
            DSL
          }
        }
        periodTotals {
          period
          NL
          DSL
        }
        totalsByFuel {
          NL
          DSL
        }
      }
    }
  `)

	req.Var("date", formattedDate(c.request.Date))
	req.Header = c.hdrs

	ctx := context.Background()
	err = c.client.Run(ctx, req, &rpt)
	if err != nil {
		log.Errorf("error running graphql client: %s", err.Error())
		return nil, err
	}
	rpt.Date = c.request.Date

	return rpt, err
}

//
// ======================== Helper Functions =============================== //
//

// extractSaleKeys helper function
func extractSaleKeys(m map[string]float64) (keys []string) {

	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		log.Errorf("input type not a map: %v", v)
		return nil
	}

	// Set expected fuel types order
	for _, ft := range model.FuelTypes {
		for _, k := range v.MapKeys() {
			if ft == k.String() {
				keys = append(keys, ft)
			}
		}
	}

	return keys
}

// sortFuelTypes function
func sortFuelTypes(fts []string) (ret []string) {
	for _, ft := range model.FuelTypes {
		for _, k := range fts {
			if ft == k {
				ret = append(ret, ft)
			}
		}
	}
	return ret
}

// formattedDate function
func formattedDate(date time.Time) string {
	return date.Format(timeLongFrmt)
}
