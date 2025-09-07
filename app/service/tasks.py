from app.exceptions import ForbiddenForUserRole
from app.models.tasks import TaskStatus
from app.repository.tasks import TaskRepository
from app.schemas.tasks import TaskCreate, TaskDetail
from app.schemas.users import User, Role


class TaskService:
    def __init__(self, task_repository: TaskRepository):
        self.task_repository: TaskRepository = task_repository

    async def create_task(self, task: TaskCreate, user: User) -> int:
        user_id = user.id
        user_role = user.role
        if user_role == Role.USER:
            if task.operator and task.operator != user_id:
                raise ForbiddenForUserRole("user with USER role cannot assign anyone other than himself")
            if task.status == TaskStatus.DONE:
                raise ForbiddenForUserRole(f"user with USER role cannot create task with `{TaskStatus.DONE}` status")
        task_dict = task.model_dump()
        watch_self = task_dict.pop("watch_self")
        task_dict["author"] = user_id
        watcher_id = user_id if watch_self else None
        task_id = await self.task_repository.create_task(task_dict, watcher_id)
        return task_id

    async def get_task(self, task_id: int) -> TaskDetail:
        task = await self.task_repository.get_task(task_id)
        task_schema = TaskDetail.model_validate(task)
        return task_schema
