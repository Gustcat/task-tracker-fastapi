from fastapi import Depends

from app.db.db import SessionDep
from app.repository.tasks import TaskRepository


async def get_task_repo(session: SessionDep):
    return TaskRepository(session)

TaskRepoDep = Depends(get_task_repo)
