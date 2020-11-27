#!/bin/bash

source secrets.sh
export HOST="0.0.0.0"
export TWITCH_CHANNEL="hennedo92"

go run .
