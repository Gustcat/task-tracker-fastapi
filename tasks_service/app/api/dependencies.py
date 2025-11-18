from typing import Annotated

from fastapi import Depends

from app.repository.dependecies import TaskRepoDep
from app.repository.tasks import TaskRepository
from app.service.tasks import TaskService
from app.сlients.dependencies import AuthGRPCClientDep
from app.сlients.auth_grpc_client import AuthGRPCClient


async def get_task_service(repo: TaskRepository = TaskRepoDep,
                           user_client: AuthGRPCClient = AuthGRPCClientDep) -> TaskService:
    return TaskService(repo, user_client)

TaskServiceDep = Annotated[TaskService, Depends(get_task_service)]
