from typing import Annotated

from fastapi import APIRouter, Depends

from app.schemas.tasks import TaskCreate
from app.schemas.users import User
from app.security import get_current_user
from app.service.tasks import TaskService
from app.api.dependencies import get_task_service


router = APIRouter()


@router.post("/")
async def create_task(task: TaskCreate,
                      task_service: Annotated[TaskService, Depends(get_task_service)],
                      user: User = Depends(get_current_user)):
    task_id = await task_service.create_task(task, user)
    return {"task_id": task_id}
