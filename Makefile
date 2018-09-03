include .env

# found yolo at: https://azer.bike/journal/a-good-makefile-for-go/

default: build

deploy: build awsPackage awsDeploy

clean:
	@rm -rf dist
	@mkdir -p dist

build: clean
	@for dir in `ls handler`; do \
		GOOS=linux go build -o dist/$$dir github.com/pulpfree/gdps-fs-dwnld/handler/$$dir; \
	done
	@cp ./config/defaults.yaml dist/
	@echo "build successful"

# watch: Run given command when code changes. e.g; make watch run="echo 'hey'"
# @yolo -i . -e vendor -e bin -e dist -c $(run)

watch:
	@yolo -i . -e vendor -e dist -c "make build"

validate:
	sam validate

run: build
	sam local start-api -n env.json

awsPackage:
	aws cloudformation package \
   --template-file template.yaml \
   --output-template-file packaged-template.yaml \
   --s3-bucket $(AWS_BUCKET_NAME) \
   --s3-prefix lambda \
   --profile $(AWS_PROFILE)

awsDeploy:
	aws cloudformation deploy \
   --template-file packaged-template.yaml \
   --stack-name $(AWS_STACK_NAME) \
   --capabilities CAPABILITY_IAM \
   --profile $(AWS_PROFILE)

describe:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(AWS_STACK_NAME)