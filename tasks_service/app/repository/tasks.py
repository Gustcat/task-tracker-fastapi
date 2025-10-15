import sqlalchemy
from sqlalchemy import select, update, delete
from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.exceptions import (
    TaskNotFoundError,
    TaskAlreadyExistsError,
    TaskWatcherAlreadyExistsError,
    TaskWatcherNotFoundError,
)
from app.models.tasks import TaskModel, TaskWatcherModel
from app.repository.utils import flush_or_raise


class TaskRepository:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def create_task(self, task_dict: dict) -> int:
        task = TaskModel(**task_dict)
        self.session.add(task)
        await flush_or_raise(self.session, TaskAlreadyExistsError, title=task.title)
        await self.session.commit()
        return task.id

    async def _fetch_task(self, query, task_id: int) -> TaskModel:
        task = (await self.session.execute(query)).scalar_one_or_none()
        if task is None:
            raise TaskNotFoundError(task_id)
        return task

    async def get_task_with_watchers(self, task_id: int) -> TaskModel:
        query = (
            select(TaskModel)
            .options(selectinload(TaskModel.watchers))
            .where(TaskModel.id == task_id)
        )
        return await self._fetch_task(query, task_id)

    async def get_basic_task(self, task_id: int) -> TaskModel:
        query = select(TaskModel).where(TaskModel.id == task_id)
        return await self._fetch_task(query, task_id)

    async def update_task(self, task: TaskModel, task_dict: dict) -> TaskModel:
        await self.session.execute(
            update(TaskModel).where(TaskModel.id == task.id).values(**task_dict)
        )
        await flush_or_raise(self.session, TaskAlreadyExistsError, title=task.title)
        await self.session.refresh(task)
        return task

    async def list_tasks(
        self, filter_dict: dict
    ) -> tuple[list[TaskModel], bool, int, int]:
        query = select(TaskModel)
        if "watcher" in filter_dict:
            query = query.join(TaskModel.watchers).where(
                TaskWatcherModel.user_id == filter_dict["watcher"]
            )
        if "author" in filter_dict:
            query = query.where(TaskModel.author == filter_dict["author"])
        if "operator" in filter_dict:
            query = query.where(TaskModel.operator == filter_dict["operator"])
        if "status" in filter_dict:
            query = query.where(TaskModel.status == filter_dict["status"])
        if "title" in filter_dict:
            query = query.where(TaskModel.title.ilike(f"%{filter_dict['title']}%"))

        limit = filter_dict["limit"]
        offset = filter_dict["offset"]
        order_field = getattr(TaskModel, filter_dict["order_by"])

        if filter_dict["is_desc"]:
            order_field = order_field.desc()
        query = query.order_by(order_field).offset(offset).limit(limit + 1)
        result = (await self.session.execute(query)).scalars().all()

        has_next = len(result) > limit
        result = result[:limit]

        return list(result), has_next, offset, limit

    async def delete_task(self, task: TaskModel) -> None:
        await self.session.delete(task)

    async def add_watcher(self, taskwatcher: TaskWatcherModel) -> None:
        self.session.add(taskwatcher)
        await flush_or_raise(
            self.session, TaskWatcherAlreadyExistsError, task_id=taskwatcher.task_id
        )

    async def remove_watcher(self, task_id: int, user_id: int) -> None:
        result = await self.session.execute(
            delete(TaskWatcherModel).where(
                TaskWatcherModel.task_id == task_id,
                TaskWatcherModel.user_id == user_id,
            )
        )
        if result.rowcount == 0:
            raise TaskWatcherNotFoundError(task_id)
