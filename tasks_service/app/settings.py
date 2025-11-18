from pathlib import Path
from typing import Annotated

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict

BASE_DIR = Path(__file__).resolve().parent
ENV_PATH = BASE_DIR.parent / ".env"


class Settings(BaseSettings):
    db_user: Annotated[str, Field(alias="POSTGRES_USER")] = "user"
    db_pass: Annotated[str, Field(alias="POSTGRES_PASSWORD")] = "password"
    db_port: Annotated[str, Field(alias="POSTGRES_PORT")] = "5432"
    db_name: Annotated[str, Field(alias="POSTGRES_DB")] = "task"
    db_echo: Annotated[bool, Field(alias="DB_ECHO")] = False
    secret_key: Annotated[str, Field(alias="ACCESS_TOKEN_SECRET")]
    algorithm: Annotated[str, Field(alias="ALGORITHM")] = "HS256"
    http_host: Annotated[str, Field(alias="HTTP_HOST")] = "localhost"
    http_port: Annotated[int, Field(alias="HTTP_PORT")] = 8081
    grpc_host: Annotated[str, Field(alias="GRPC_HOST")] = "localhost"
    grpc_port: Annotated[int, Field(alias="GRPC_PORT")] = 50051

    model_config = SettingsConfigDict(
        env_file=ENV_PATH, env_file_encoding="utf-8", extra="ignore"
    )

    @property
    def REAL_DATABASE_URL(self):
        return f"postgresql+asyncpg://{settings.db_user}:{settings.db_pass}@localhost:{settings.db_port}/{settings.db_name}"

    @property
    def GRPC_ADDRESS(self):
        return f"{settings.grpc_host}:{settings.grpc_port}"


settings = Settings()
