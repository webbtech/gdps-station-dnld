# Create SSM Params

### Create Parameters For Development

* alias is: alias/testPulpfree
* path is: /test/gdps-fs-dwnld

``` bash
$ aws ssm put-parameter --name /test/gdps-fs-dwnld/S3Bucket \
  --value gdps-fs-dwnld --type String --overwrite

$ aws ssm put-parameter --name /test/gdps-fs-dwnld/CognitoClientID \
  --value us-east-1_gsB59wfzW --type String --overwrite

$ aws ssm put-parameter --name /test/gdps-fs-dwnld/CognitoPoolID \
  --value 2084ukslsc831pt202t2dudt7c --type String --overwrite

$ aws ssm put-parameter --name /test/gdps-fs-dwnld/CognitoRegion \
  --value us-east-1 --type String --overwrite

# delete a parameter
$ aws ssm delete-parameter --name /test/gdps-fs-dwnld/DBName

# fetch params by path
$ aws ssm get-parameters-by-path --path /test/gdps-fs-dwnld
```

### Production Parameters
* alias is: alias/GalesProd
* path is: /prod/gdps-fs-dwnld

``` bash
$ aws ssm put-parameter --name /prod/gdps-fs-dwnld/S3Bucket \
  --value gdps-fs-dwnld --type String --overwrite

$ aws ssm put-parameter --name /prod/gdps-fs-dwnld/CognitoClientID \
  --value us-east-1_gsB59wfzW --type String --overwrite

$ aws ssm put-parameter --name /prod/gdps-fs-dwnld/CognitoPoolID \
  --value 2084ukslsc831pt202t2dudt7c --type String --overwrite

$ aws ssm put-parameter --name /prod/gdps-fs-dwnld/CognitoRegion \
  --value us-east-1 --type String --overwrite

# This won't work as ssm tries to lookup the url, see next one below for method that works
$ aws ssm put-parameter --name /prod/gdps-fs-dwnld/GraphqlURI \
  --value "https://fqhx4fm4d8.execute-api.ca-central-1.amazonaws.com/Prod/graphql" --type String --overwrite

# Use the cli-input-json to use a url as value
$ aws ssm put-parameter --cli-input-json '{
  "Name": "/prod/gdps-fs-dwnld/GraphqlURI",
  "Value": "https://fqhx4fm4d8.execute-api.ca-central-1.amazonaws.com/Prod/graphql",
  "Type": "String",
  "Description": "url",
  "Overwrite": true
}'

```