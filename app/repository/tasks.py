from fastapi import HTTPException
from asyncpg.exceptions import UniqueViolationError
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession

from app.models.tasks import TaskModel, TaskWatcherModel


class TaskRepository:
    def __init__(self, session: AsyncSession):
        self.session = session
        self.model = TaskModel

    async def create_task(self, task_data: dict,  watcher_id: int | None = None) -> int:
        task = self.model(**task_data)
        self.session.add(task)
        if watcher_id is not None:
            task.watchers.append(TaskWatcherModel(user_id=watcher_id))
        try:
            await self.session.commit()
        except IntegrityError as e:
            await self.session.rollback()
            if getattr(e.orig, 'pgcode', None) == '23505':
                raise HTTPException(status_code=400, detail="task with such title already exists")
            raise HTTPException(status_code=500, detail="Integrity error")
        return task.id

    async def get_task(self, task_id: int) -> TaskModel:
        pass