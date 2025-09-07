from fastapi import FastAPI
import uvicorn

from app.api.tasks import router
from app.errors_handlers import register_errors_handlers
from app.settings import settings

app = FastAPI()

register_errors_handlers(app)

app.include_router(router, prefix='/tasks', tags=["tasks"])

if __name__ == "__main__":
    uvicorn.run(app, host=settings.http_host, port=settings.http_port)


