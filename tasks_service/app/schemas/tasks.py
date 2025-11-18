from datetime import date, datetime
from typing import Literal, Annotated

from annotated_types import MinLen, MaxLen
from fastapi import HTTPException
from pydantic import BaseModel, field_validator, ConfigDict, model_validator

from app.models.tasks import TaskStatus

title = Annotated[str, MinLen(2), MaxLen(250)]


class TaskChangeSchema(BaseModel):
    description: str | None = None
    operator: int | None = None
    target_date: date | None = None
    status: TaskStatus | None = None

    @field_validator("target_date")
    def validate_target_date(cls, value):
        if value and date.today() > value:
            raise HTTPException(
                status_code=422, detail="due date should not earlier then now"
            )
        return value


class TaskCreateSchema(TaskChangeSchema):
    title: title

    model_config = ConfigDict(extra="forbid")

    @model_validator(mode="after")
    def set_default_status(cls, values):
        if not values.status:
            values.status = TaskStatus.TODO if values.operator else TaskStatus.NEW
        return values


class TaskUpdateSchema(TaskChangeSchema):
    title: title = None
    status: TaskStatus | None = None

    @field_validator("title", "description", "status", mode="before")
    def forbid_null_for_some_fields(cls, value, info):
        if value is None:
            raise ValueError(
                f"Field '{info.field_name}' cannot be null. "
                f"Omit it if you don't want to update this field."
            )
        return value


class TaskListSchema(BaseModel):
    id: int
    author: int
    author_deleted: bool
    title: title
    description: str | None = None
    status: TaskStatus = TaskStatus.NEW
    operator: int | None = None
    operator_deleted: bool

    model_config = ConfigDict(from_attributes=True)


class TaskDetailSchema(TaskListSchema):
    watchers: list[int]
    target_date: date | None = None
    completed_at: datetime | None = None
    created_at: datetime = None
    updated_at: datetime = None


class TaskFilter(BaseModel):
    author: int | None = None
    operator: int | None = None
    watcher: int | None = None
    title: str | None = None
    status: TaskStatus | None = None
    order_by: Literal["created_at", "updated_at", "completed_at"] = "created_at"
    is_desc: bool = True
