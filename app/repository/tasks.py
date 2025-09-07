from sqlalchemy import select
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.exceptions import TaskNotFoundError, TaskAlreadyExistsError
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
                raise TaskAlreadyExistsError(task.title)
            raise
        return task.id

    async def get_task(self, task_id: int) -> TaskModel:
        query = (
            select(TaskModel)
            .options(selectinload(TaskModel.watchers))
            .where(TaskModel.id == task_id)
        )
        task = (await self.session.execute(query)).scalar_one_or_none()
        if task is None:
            raise TaskNotFoundError(task_id)
        return task
