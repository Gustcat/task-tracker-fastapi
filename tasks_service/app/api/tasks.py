from typing import Annotated

from fastapi import APIRouter, Depends, status, Path, Query

from app.schemas.pagination import PaginatedResponse, PaginationParams
from app.schemas.tasks import (
    TaskCreateSchema,
    TaskDetailSchema,
    TaskFilter,
    TaskUpdateSchema,
    TaskListSchema
)
from app.schemas.users import User
from app.security import get_current_user
from app.api.dependencies import TaskServiceDep

router = APIRouter()


@router.post("/", status_code=status.HTTP_201_CREATED)
async def create_task(
    task: TaskCreateSchema,
    task_service: TaskServiceDep,
    user: User = Depends(get_current_user),
):
    task_id = await task_service.create_task(task, user)
    return {"task_id": task_id}


@router.get("/", response_model=PaginatedResponse)
async def list_tasks(
    pagination: Annotated[PaginationParams, Depends()],
    task_filter: Annotated[TaskFilter, Query()],
    task_service: TaskServiceDep,
    _: User = Depends(get_current_user),
):

    paginated_response = await task_service.list_tasks(task_filter, pagination)
    return paginated_response


@router.get("/{task_id}/", response_model=TaskDetailSchema)
async def get_task(
    task_id: Annotated[int, Path(ge=1)],
    task_service: TaskServiceDep,
    _: User = Depends(get_current_user),
):
    task = await task_service.get_task(task_id)
    return task


@router.patch("/{task_id}/", response_model=TaskListSchema)
async def update_task(
    task_id: Annotated[int, Path(ge=1)],
    task_service: TaskServiceDep,
    task_update_schema: TaskUpdateSchema,
    user: User = Depends(get_current_user),
):
    task = await task_service.update_task(task_id, user, task_update_schema)
    return task


@router.delete("/{task_id}/", status_code=status.HTTP_204_NO_CONTENT)
async def delete_task(
    task_id: Annotated[int, Path(ge=1)],
    task_service: TaskServiceDep,
    user: User = Depends(get_current_user),
):
    await task_service.delete_task(task_id, user)
    return


@router.post("/{task_id}/watchers/me/", status_code=status.HTTP_204_NO_CONTENT)
async def add_self_watcher(
    task_id: Annotated[int, Path(ge=1)],
    task_service: TaskServiceDep,
    user: User = Depends(get_current_user),
):
    await task_service.add_watcher(task_id, user.id)
    return


@router.delete("/{task_id}/watchers/me/", status_code=status.HTTP_204_NO_CONTENT)
async def remove_self_watcher(
    task_id: Annotated[int, Path(ge=1)],
    task_service: TaskServiceDep,
    user: User = Depends(get_current_user),
):
    await task_service.remove_watcher(task_id, user.id)
    return
