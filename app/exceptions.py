class AppError(Exception):
    pass


class TaskNotFoundError(AppError):
    def __init__(self, task_id: int):
        super().__init__(f"Task with {task_id=} not found")


class TaskAlreadyExistsError(AppError):
    def __init__(self, title: str):
        super().__init__(f"Task with {title=} already exists")


class ForbiddenForUserRole(AppError):
    def __init__(self, msg_end: str):
        super().__init__(f"User with USER role cannot  {msg_end}")
