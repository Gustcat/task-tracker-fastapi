from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    db_user: str = Field(default="user", alias="POSTGRES_USER")
    db_pass: str = Field(default="password", alias="POSTGRES_PASSWORD")
    db_port: str = Field(default="5432", alias="POSTGRES_PORT")
    db_name: str = Field(default="task", alias="POSTGRES_DB")
    secret_key: str = Field(..., alias="ACCESS_TOKEN_SECRET")
    algorithm: str = Field(default="HS256", alias="ALGORITHM")
    http_host: str = Field(default="localhost", alias="HTTP_HOST")
    http_port: int = Field(default=8081, alias="HTTP_PORT")

    model_config = SettingsConfigDict(
        env_file="../.env", env_file_encoding="utf-8", extra="ignore"
    )


settings = Settings()
print(settings)

REAL_DATABASE_URL = f"postgresql+asyncpg://{settings.db_user}:{settings.db_pass}@localhost:{settings.db_port}/{settings.db_name}"
