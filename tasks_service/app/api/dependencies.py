from app.repository.dependecies import TaskRepoDep
from app.repository.tasks import TaskRepository
from app.service.tasks import TaskService
from app.сlients.dependencies import UserClientDep
from app.сlients.user_client import UserClient


async def get_task_service(repo: TaskRepository = TaskRepoDep,
                           user_client: UserClient = UserClientDep) -> TaskService:
    return TaskService(repo, user_client)
