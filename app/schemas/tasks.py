from datetime import date, datetime
from typing import Literal

from fastapi import HTTPException
from pydantic import BaseModel, Field, field_validator, ConfigDict

from app.models.tasks import TaskStatus


class TaskCreateSchema(BaseModel):
    title: str = Field(min_length=2, max_length=250)
    description: str | None = None
    status: TaskStatus = TaskStatus.NEW
    watch_self: bool = False
    operator: int | None = None
    target_date: date | None = None

    @field_validator("target_date")
    def validate_target_date(cls, value):
        if date.today() > value:
            raise HTTPException(
                status_code=422, detail="due date should not earlier then now"
            )
        return value


class TaskListSchema(BaseModel):
    id: int
    author: int
    author_deleted: bool
    title: str = Field(min_length=2, max_length=250)
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


class TaskUpdateSchema(BaseModel):
    title: str | None = Field(min_length=2, max_length=250, default=None)
    description: str | None = None
    status: TaskStatus | None = None
    operator: int | None = None
    target_date: date | None = None
    watch_self: bool | None = None


class BaseFilter(BaseModel):
    limit: int = Field(25, gt=0)
    offset: int = Field(0, ge=0)
    order_by: Literal["created_at", "updated_at"] = "created_at"

    model_config = ConfigDict(extra="forbid")


class TaskFilter(BaseFilter):
    author: int | None = None
    operator: int | None = None
    watcher: int | None = None
    title: str | None = None
    status: TaskStatus | None = None
    order_by: Literal["created_at", "updated_at", "completed_at"] = "created_at"
    is_desc: bool = True


class PaginatedResponse(BaseModel):
    items: list[TaskListSchema]
    has_next: bool
    limit: int = Field(gt=0)
    offset: int = Field(ge=0)
