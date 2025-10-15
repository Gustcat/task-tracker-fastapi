#!/bin/bash
source prod.env

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${DSN}" up -v