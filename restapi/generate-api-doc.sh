#!/usr/bin/env bash
#
# This script generates a fresh swagger.yaml and swagger.json
# Run only if you have changed the declarative comments in server.go or jobs.go
#
# @author rouven.himmelstein@cgm.com
# @date 19.08.2020
#
##############

go get -u github.com/swaggo/swag/cmd/swag
swag init -g server.go
