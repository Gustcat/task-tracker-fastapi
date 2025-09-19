from datetime import datetime

from app.exceptions import ForbiddenForUserRole, InvalidTaskStateError
from app.models.tasks import TaskStatus, TaskModel, TaskWatcherModel
from app.repository.tasks import TaskRepository
from app.schemas.tasks import (
    TaskCreateSchema,
    TaskDetailSchema,
    TaskFilter,
    TaskListSchema,
    PaginatedResponse,
    TaskUpdateSchema,
)
from app.schemas.users import User, Role


def validate_input_task(user: User, old_task: TaskModel | None, task_dict: dict):
    user_id = user.id
    user_role = user.role

    is_new_operator_exists = "operator" in task_dict
    new_operator = task_dict.get("operator")
    new_title = task_dict.get("title")
    new_status = task_dict.get("status")

    if user_role == Role.USER:
        if is_new_operator_exists and new_operator != user_id:
            raise ForbiddenForUserRole("assign anyone other than himself")

        if old_task is None:
            if new_status == TaskStatus.DONE:
                raise ForbiddenForUserRole("create task with DONE status")
        else:
            if user_id != old_task.author and user_id != old_task.operator:
                raise ForbiddenForUserRole("update other people's tasks")
            if (
                new_title is not None
                and user_id != old_task.author
                and new_title != old_task.title
            ):
                raise ForbiddenForUserRole("change title if not author")

    if new_status == TaskStatus.DONE:
        task_dict["completed_at"] = datetime.utcnow
    elif (
        old_task
        and old_task.status == TaskStatus.DONE
        and new_status
        and new_status != TaskStatus.DONE
    ):
        task_dict["completed_at"] = None

    actual_operator = (
        new_operator
        if is_new_operator_exists
        else (old_task.operator if old_task else None)
    )
    actual_status = (
        new_status
        if new_status is not None
        else (old_task.status if old_task else None)
    )

    if old_task and actual_status == TaskStatus.NEW:
        raise InvalidTaskStateError("Task cannot be moved back to NEW")

    if (
        actual_status in {TaskStatus.DONE, TaskStatus.IN_PROGRESS}
        and not actual_operator
    ):
        raise InvalidTaskStateError(
            "Task cannot be DONE or IN_PROGRESS without operator"
        )
    return


class TaskService:
    def __init__(self, task_repository: TaskRepository):
        self.task_repository: TaskRepository = task_repository

    async def create_task(self, task_schema: TaskCreateSchema, user: User) -> int:
        task_dict = task_schema.model_dump()
        validate_input_task(user, None, task_dict)
        task_dict["author"] = user.id
        task_id = await self.task_repository.create_task(task_dict)
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
            exclude_unset=True,
        )
        tasks, has_next, offset, limit = await self.task_repository.list_tasks(
            filter_dict
        )
        return PaginatedResponse(
            items=tasks, has_next=has_next, offset=offset, limit=limit
        )

    async def update_task(
        self, task_id: int, user: User, task_update_schema: TaskUpdateSchema
    ) -> TaskListSchema:
        async with self.task_repository.session.begin():
            task = await self.task_repository.get_basic_task(task_id)
            task_dict = task_update_schema.model_dump(exclude_unset=True)
            validate_input_task(user, task, task_dict)
            updated_task = await self.task_repository.update_task(task, task_dict)
            task_schema = TaskListSchema.model_validate(updated_task)
            return task_schema

    async def delete_task(self, task_id: int, user: User):
        async with self.task_repository.session.begin():
            task = await self.task_repository.get_basic_task(task_id)
            if user.role == Role.USER and user.id != task.author:
                raise ForbiddenForUserRole(f"delete other people's tasks")
            await self.task_repository.delete_task(task)

    async def add_watcher(self, task_id: int, user_id: int):
        async with self.task_repository.session.begin():
            await self.task_repository.get_basic_task(task_id)
            task_watcher = TaskWatcherModel(task_id=task_id, user_id=user_id)
            await self.task_repository.add_watcher(task_watcher)

    async def remove_watcher(self, task_id: int, user_id: int):
        async with self.task_repository.session.begin():
            await self.task_repository.get_basic_task(task_id)
            await self.task_repository.remove_watcher(task_id, user_id)
