from datetime import date, datetime

from fastapi import HTTPException
from pydantic import BaseModel, Field, field_validator

from app.models.tasks import TaskStatus


class TaskCreate(BaseModel):
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
    description: str = None
    status: TaskStatus = TaskStatus.NEW
    operator: int = None
    operator_deleted: bool

    class Config:
        from_attributes = True

    def to_response(self):
        data = self.model_dump()
        if self.operator is None:
            data.pop("operator_deleted", None)
        return data


class TaskDetail(TaskListSchema):
    watchers: list[int]
    target_date: date = None
    completed_at: datetime = None
    created_at: datetime = None
    updated_at: datetime = None


class TaskUpdate(BaseModel):
    title: str | None = Field(min_length=2, max_length=250, default=None)
    description: str | None = None
    status: TaskStatus | None = None
    operator: int | None = None
    target_date: date | None = None
    watch_self: bool | None = None
