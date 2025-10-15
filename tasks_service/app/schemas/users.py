from enum import IntEnum

from pydantic import BaseModel


class Role(IntEnum):
    USER = 1
    ADMIN = 2


class User(BaseModel):
    id: int
    name: str
    role: Role
