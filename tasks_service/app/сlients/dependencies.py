from fastapi import Depends

from app.settings import GRPC_ADDRESS
from app.сlients.user_client import UserClient


async def get_task_client() -> UserClient:
    return UserClient(GRPC_ADDRESS)

UserClientDep = Depends(get_task_client)
