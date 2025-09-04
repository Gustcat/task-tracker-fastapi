from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession
from typing_extensions import Annotated

from app.repository.tasks import TaskRepository
from app.service.tasks import TaskService
from db.db import get_async_session


def get_task_service(session: Annotated[AsyncSession, Depends(get_async_session)]):
    return TaskService(TaskRepository(session))
