from envparse import Env
from fastapi import FastAPI
from fastapi.routing import APIRouter
import uvicorn

from app.api.tasks import router
from app.settings import settings

app = FastAPI()

app.include_router(router, prefix='/tasks', tags=["tasks"])

if __name__ == "__main__":
    uvicorn.run(app, host=settings.http_host, port=settings.http_port)


