package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	config
	DefaultsFilePath string
}

// defaults struct
type defaults struct {
	AWSRegion  string `yaml:"AWSRegion"`
	S3Bucket   string `yaml:"S3Bucket"`
	GraphqlURI string `yaml:"GraphqlURI"`
	Stage      string `yaml:"Stage"`
}

type config struct {
	AWSRegion  string
	S3Bucket   string
	GraphqlURI string
	Stage      StageEnvironment
}

// Dynamo struct
type Dynamo struct {
	APIVersion string `yaml:"APIVersion"`
	Region     string `yaml:"Region"`
}

// StageEnvironment string
type StageEnvironment string

// StageEnvironment type constants
const (
	DevEnv   StageEnvironment = "dev"
	StageEnv StageEnvironment = "stage"
	TestEnv  StageEnvironment = "test"
	ProdEnv  StageEnvironment = "prod"
)

const defaultFileName = "defaults.yaml"

var (
	defs = &defaults{}
)

// Load method
func (c *Config) Load() (err error) {

	if err = c.setDefaults(); err != nil {
		return err
	}

	if err = c.setEnvVars(); err != nil {
		return err
	}

	c.setFinal()

	return err
}

// GetStageEnv method
func (c *Config) GetStageEnv() StageEnvironment {
	return c.Stage
}

// this must be called first in c.Load
func (c *Config) setDefaults() (err error) {

	if c.DefaultsFilePath == "" {
		dir, _ := os.Getwd()
		c.DefaultsFilePath = path.Join(dir, defaultFileName)
	}

	file, err := ioutil.ReadFile(c.DefaultsFilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(file), &defs)
	if err != nil {
		return err
	}
	err = c.validateStage()

	return err
}

// validateStage method to validate Stage value
func (c *Config) validateStage() (err error) {

	validEnv := true

	switch defs.Stage {
	case "dev":
	case "development":
		c.Stage = DevEnv
	case "stage":
		c.Stage = StageEnv
	case "test":
		c.Stage = TestEnv
	case "prod":
		c.Stage = ProdEnv
	case "production":
		c.Stage = ProdEnv
	default:
		validEnv = false
	}

	if !validEnv {
		return errors.New(fmt.Sprintf("Invalid StageEnvironment requested: %s", defs.Stage))
	}

	return err
}

// sets any environment variables that match the default struct fields
func (c *Config) setEnvVars() (err error) {

	vals := reflect.Indirect(reflect.ValueOf(defs))
	for i := 0; i < vals.NumField(); i++ {
		nm := vals.Type().Field(i).Name
		if e := os.Getenv(nm); e != "" {
			vals.Field(i).SetString(e)
		}
		// If field is Stage, validate and return error if required
		if nm == "Stage" {
			err = c.validateStage()
			if err != nil {
				return err
			}
		}
	}

	return err
}

// Copies required fields from the defaults to the Config struct
func (c *Config) setFinal() {
	c.AWSRegion = defs.AWSRegion
	c.GraphqlURI = defs.GraphqlURI
	c.S3Bucket = defs.S3Bucket
}
