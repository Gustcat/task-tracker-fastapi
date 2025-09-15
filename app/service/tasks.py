from app.exceptions import ForbiddenForUserRole
from app.models.tasks import TaskStatus, TaskModel
from app.repository.tasks import TaskRepository
from app.schemas.tasks import (
    TaskCreateSchema,
    TaskDetailSchema,
    TaskFilter,
    TaskListSchema,
    PaginatedResponse,
)
from app.schemas.users import User, Role


class TaskService:
    def __init__(self, task_repository: TaskRepository):
        self.task_repository: TaskRepository = task_repository

    async def create_task(self, task_schema: TaskCreateSchema, user: User) -> int:
        user_id = user.id
        user_role = user.role
        if user_role == Role.USER:
            if task_schema.operator and task_schema.operator != user_id:
                raise ForbiddenForUserRole("assign anyone other than himself")
            if task_schema.status == TaskStatus.DONE:
                raise ForbiddenForUserRole(
                    f"create task with `{TaskStatus.DONE}` status"
                )
        task_dict = task_schema.model_dump()
        watch_self = task_dict.pop("watch_self")
        task_dict["author"] = user_id
        task = TaskModel(**task_dict)
        watcher_id = user_id if watch_self else None
        task_id = await self.task_repository.create_task(task, watcher_id)
        return task_id

    async def get_task(self, task_id: int) -> TaskDetailSchema:
        task = await self.task_repository.get_task_with_watchers(task_id)
        task_data = task.__dict__.copy()
        task_data["watchers"] = [w.user_id for w in task.watchers]
        # добавить batch-запрос к user-сервису для всех user_id
        task_schema = TaskDetailSchema.model_validate(task_data)
        return task_schema

    async def list_tasks(self, task_filter: TaskFilter) -> PaginatedResponse:
        filter_dict = task_filter.model_dump(
            exclude_none=True,
        )
        tasks, has_next, offset, limit = await self.task_repository.list_tasks(
            filter_dict
        )
        return PaginatedResponse(
            items=tasks, has_next=has_next, offset=offset, limit=limit
        )

    async def delete_task(self, task_id: int, user: User):
        async with self.task_repository.session.begin():
            task = await self.task_repository.get_basic_task(task_id)
            if user.role == Role.USER and user.id != task.author:
                raise ForbiddenForUserRole(f"delete other people's tasks")
            await self.task_repository.delete_task(task)
