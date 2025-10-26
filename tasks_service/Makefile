include .env

DB_URL="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}"
ALEMBIC_DIR = alembic

PROTO_DIR_USER=../auth/api/user_v1
PROTO_DIR_PROTOGEN=../auth/vendor.protogen
VENV_PACKAGE_DIR=venv/Lib/site-packages

OUT_DIR_USER=app/proto/user
OUT_DIR_AUTH=app/proto/auth

init_alembic:
	echo ${DB_URL}
	@if [ ! -d ${ALEMBIC_DIR} ]; then \
		alembic init ${ALEMBIC_DIR}; \
	fi
	@sed -i "s|^sqlalchemy.url = .*|sqlalchemy.url = ${DB_URL}|" alembic.ini

generate_validate:
	python -m grpc_tools.protoc \
	-I${PROTO_DIR_PROTOGEN} \
	--python_out=${VENV_PACKAGE_DIR} \
	--grpc_python_out=${VENV_PACKAGE_DIR} \
	${PROTO_DIR_PROTOGEN}/validate/validate.proto

generate_user_api:
	mkdir -p ${OUT_DIR_USER}
	python -m grpc_tools.protoc \
	-I${PROTO_DIR_USER} \
	-I${PROTO_DIR_PROTOGEN} \
	--python_out=${OUT_DIR_USER} \
	--grpc_python_out=${OUT_DIR_USER} \
	--pyi_out=${OUT_DIR_USER} \
	${PROTO_DIR_USER}/*.proto