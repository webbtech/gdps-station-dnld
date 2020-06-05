# Create SSM Params

* cognito pool id: ca-central-1_lolwfYIAr
* cognito client id: 5n63nd473pv7ne2qskv30gkcbh

## Create Parameters For Development

* alias is: alias/testPulpfree
* path is: /test/gdps-fs-dwnld

``` bash
aws ssm put-parameter --name /test/gdps-fs-dwnld/S3Bucket \
  --value gdps-fs-dwnld --type String --overwrite

aws ssm put-parameter --name /test/gdps-fs-dwnld/CognitoClientID \
  --value ca-central-1_lolwfYIAr --type String --overwrite

aws ssm put-parameter --name /test/gdps-fs-dwnld/CognitoPoolID \
  --value 5n63nd473pv7ne2qskv30gkcbh --type String --overwrite

aws ssm put-parameter --name /test/gdps-fs-dwnld/CognitoRegion \
  --value ca-central-1 --type String --overwrite

# delete a parameter
aws ssm delete-parameter --name /test/gdps-fs-dwnld/DBName

# fetch params by path
aws ssm get-parameters-by-path --path /test/gdps-fs-dwnld
```

## Production Parameters

* alias is: alias/GalesProd
* path is: /prod/gdps-fs-dwnld

``` bash
aws ssm put-parameter --name /prod/gdps-fs-dwnld/S3Bucket \
  --value gdps-fs-dwnld --type String --overwrite

aws ssm put-parameter --name /prod/gdps-fs-dwnld/CognitoClientID \
  --value ca-central-1_lolwfYIAr --type String --overwrite

aws ssm put-parameter --name /prod/gdps-fs-dwnld/CognitoPoolID \
  --value 5n63nd473pv7ne2qskv30gkcbh --type String --overwrite

aws ssm put-parameter --name /prod/gdps-fs-dwnld/CognitoRegion \
  --value ca-central-1 --type String --overwrite

# This won't work as ssm tries to lookup the url, see next one below for method that works
aws ssm put-parameter --name /prod/gdps-fs-dwnld/GraphqlURI \
  --value "https://api-prod.gdps.pfapi.io/graphql" --type String --overwrite

# Use the cli-input-json to use a url as value
aws ssm put-parameter --cli-input-json '{
  "Name": "/prod/gdps-fs-dwnld/GraphqlURI",
  "Value": "https://api-prod.gdps.pfapi.io/graphql",
  "Type": "String",
  "Description": "url",
  "Overwrite": true
}'

# fetch params by path
aws ssm get-parameters-by-path --path /prod/gdps-fs-dwnld

```
