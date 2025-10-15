#!/bin/bash
source local.env

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${DSN}" up -v