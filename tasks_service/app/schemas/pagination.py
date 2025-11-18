from typing import Annotated, TypeVar, Generic

from annotated_types import Gt, Ge, Le
from pydantic import BaseModel, ConfigDict


class PaginationParams(BaseModel):
    limit: Annotated[int, Gt(0), Le(100)] = 10
    offset: Annotated[int, Ge(0)] = 0

    model_config = ConfigDict(extra="forbid")


T = TypeVar("T", bound=BaseModel)


class PaginatedResponse(BaseModel, Generic[T]):
    items: list[T]
    has_next: bool
    limit: Annotated[int, Gt(0)]
    offset: Annotated[int, Ge(0)]

    model_config = ConfigDict(extra="forbid")
