ifneq ("$(wildcard .preset)","")
$(info Loading preset .preset) 
include .preset
else
$(info .preset not found)
endif


ifeq (${stage},)
ifeq ($(STAGE),)
	stage := $(shell read -p "Enter stage: " input; echo $$input)
else
	stage := $(STAGE)
endif
endif
account_id := $(shell aws sts get-caller-identity --output json | jq '.Account')
start_time := $(shell date +%s)
project := "goserverless"
app_bucket := $(shell echo "${account_id}-${stage}-${project}-bucket")
deployment_bucket := $(shell echo "${account_id}-${stage}-${project}-deployment-bucket")
table := $(shell echo "${stage}-${project}-table")


GOBINARY := go
GOARCH := amd64
GOOS := linux
GOBUILDFLAGS := -ldflags="-s -w"
GOBUILD_CMD := 
ifneq (${GOARCH},)
GOBUILD_CMD := $(GOBUILD_CMD) GOARCH=$(GOARCH)
endif

ifneq (${GOOS},)
GOBUILD_CMD := $(GOBUILD_CMD) GOOS=$(GOOS)
endif

GOBUILD_CMD := $(GOBUILD_CMD) $(GOBINARY) build

ifneq (${GOBUILDFLAGS},)
GOBUILD_CMD := $(GOBUILD_CMD) $(GOBUILDFLAGS)
endif

define EMPTY_BUCKET_SCRIPT
import boto3
import sys
bucket_names = sys.argv[1:]
s3 = boto3.resource('s3')
for bucket_name in bucket_names:
	try:
		bucket = s3.Bucket(bucket_name)
		bucket.objects.delete()
		bucket.object_versions.delete()
	except Exception as e:
		print("[", bucket_name, "]", e)
endef

export EMPTY_BUCKET_SCRIPT

handlers := $(wildcard ./handler/*)
build: $(handlers)
	for handler in ${handlers}; \
	do \
		$(GOBUILD_CMD) -o .build/$$handler $$handler; \
	done;

define MERGE_YAML_SCRIPT
import yaml
import sys

from yamlinclude import YamlIncludeConstructor

YamlIncludeConstructor.add_to_loader_class(loader_class=yaml.FullLoader, base_dir='.')
yaml_file = sys.argv[1]
with open(yaml_file) as f:
    data = yaml.load(f, Loader=yaml.FullLoader)

print("# THIS FILE IS AUTOGENERATED. DO NOT EDIT THIS FILE DIRECTLY.")
print("# All manual changes on this file will be overwritten.")
print("# You should edit serverless.raw.yml instead.")
print("# Run make merge-yaml to update this file.")
yaml.dump(data, sys.stdout)

endef
export MERGE_YAML_SCRIPT

merge-yaml:
	@ python3 -c "$$MERGE_YAML_SCRIPT" serverless.raw.yml > serverless.yml


deploy deploy-function cleanup: %: real-%
	@ echo "================ Make $@ action done in $$(($$(date +%s) - ${start_time})) seconds"

real-cleanup:
	@ echo "Cleaning stage ${stage}..."
	@ echo "Empty bucket ${deployment_bucket}"
	@ python3 -c "$$EMPTY_BUCKET_SCRIPT" ${deployment_bucket}
	@ sls remove --stage ${stage} --verbose || true
	@ aws cloudformation delete-stack --stack-name ${project}-${stage}
	@ aws s3api delete-bucket --bucket ${deployment_bucket} || true

real-deploy:
	@ export $(grep -v '^#' .env | xargs) > /dev/null
	@ echo "Import env variables from .env"
	@ echo "Deploying serverless stage ${stage}"
	@ sls deploy --stage ${stage} --verbose

real-deploy-function: need-function
	@ echo "Deploying function ${function} stage ${stage}"
	@ sls deploy function -f ${function} --stage ${stage} --verbose


need-function:
ifeq (${function},)
ifeq ($(FUNCTION),)
	$(eval function := $(shell read -p "Enter function: " input; echo $$input))
else
	$(eval function := $(FUNCTION))
endif
endif