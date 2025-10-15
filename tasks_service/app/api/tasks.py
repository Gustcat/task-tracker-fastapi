from typing import Annotated

from fastapi import APIRouter, Depends, Query, status

from app.schemas.tasks import (
    TaskCreateSchema,
    TaskDetailSchema,
    TaskFilter,
    PaginatedResponse,
    TaskUpdateSchema,
    TaskListSchema,
)
from app.schemas.users import User
from app.security import get_current_user
from app.service.tasks import TaskService
from app.api.dependencies import get_task_service


router = APIRouter()


@router.post("/", status_code=status.HTTP_201_CREATED)
async def create_task(
    task: TaskCreateSchema,
    task_service: Annotated[TaskService, Depends(get_task_service)],
    user: User = Depends(get_current_user),
):
    task_id = await task_service.create_task(task, user)
    return {"task_id": task_id}


@router.get("/", response_model=PaginatedResponse)
async def list_tasks(
    task_filter: Annotated[TaskFilter, Query()],
    task_service: Annotated[TaskService, Depends(get_task_service)],
    _: User = Depends(get_current_user),
):
    tasks = await task_service.list_tasks(task_filter)
    return tasks


@router.get("/{task_id}/", response_model=TaskDetailSchema)
async def get_task(
    task_id: int,
    task_service: Annotated[TaskService, Depends(get_task_service)],
    _: User = Depends(get_current_user),
):
    task = await task_service.get_task(task_id)
    return task


@router.patch("/{task_id}/", response_model=TaskListSchema)
async def get_task(
    task_id: int,
    task_service: Annotated[TaskService, Depends(get_task_service)],
    task_update_schema: TaskUpdateSchema,
    user: User = Depends(get_current_user),
):
    task = await task_service.update_task(task_id, user, task_update_schema)
    return task


@router.delete("/{task_id}/", status_code=status.HTTP_204_NO_CONTENT)
async def delete_task(
    task_id: int,
    task_service: Annotated[TaskService, Depends(get_task_service)],
    user: User = Depends(get_current_user),
):
    await task_service.delete_task(task_id, user)
    return


@router.post("/{task_id}/watchers/me/", status_code=status.HTTP_201_CREATED)
async def add_self_watcher(
    task_id: int,
    task_service: Annotated[TaskService, Depends(get_task_service)],
    user: User = Depends(get_current_user),
):
    await task_service.add_watcher(task_id, user.id)
    return


@router.delete("/{task_id}/watchers/me/", status_code=status.HTTP_204_NO_CONTENT)
async def remove_self_watcher(
    task_id: int,
    task_service: Annotated[TaskService, Depends(get_task_service)],
    user: User = Depends(get_current_user),
):
    await task_service.remove_watcher(task_id, user.id)
    return
