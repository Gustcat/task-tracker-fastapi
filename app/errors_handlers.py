from http import HTTPStatus
from typing import Type

from fastapi import FastAPI, Request
from starlette.responses import JSONResponse

from app.exceptions import AppError, TaskNotFoundError, TaskAlreadyExistsError, ForbiddenForUserRole

APP_ERROR_STATUS_MAP: dict[Type[AppError], HTTPStatus] = {
    TaskNotFoundError: HTTPStatus.NOT_FOUND,
    TaskAlreadyExistsError: HTTPStatus.BAD_REQUEST,
    ForbiddenForUserRole: HTTPStatus.FORBIDDEN
}


def register_errors_handlers(app: FastAPI) -> None:

    @app.exception_handler(AppError)
    def handle_app_error(
            request: Request,
            exc: AppError
    ):
        status_code = APP_ERROR_STATUS_MAP.get(type(exc), 500)
        return JSONResponse(
            status_code=status_code,
            content={"error": str(exc)}
        )
