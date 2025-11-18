from fastapi import Depends

from app.settings import settings
from app.сlients.auth_grpc_client import AuthGRPCClient


async def get_task_client() -> AuthGRPCClient:
    return AuthGRPCClient(settings.GRPC_ADDRESS)

AuthGRPCClientDep = Depends(get_task_client)
