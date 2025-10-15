from typing import Any, Type

from sqlalchemy.exc import IntegrityError
from sqlalchemy.ext.asyncio import AsyncSession

UNIQUENESS_VIOLATION_PGCODE = "23505"


async def flush_or_raise(
    session: AsyncSession,
    error_cls: Type[Exception] | None = None,
    **kwargs: Any,
) -> None:
    """
    Безопасный flush с rollback при ошибках.

    :param session: SQLAlchemy AsyncSession
    :param error_cls: исключение, которое нужно выбросить при нарушении уникальности
    :param kwargs: аргументы для передачи в исключение error_cls
    """
    try:
        await session.flush()
    except IntegrityError as e:
        await session.rollback()
        if (
            error_cls is not None
            and getattr(e.orig, "pgcode", None) == UNIQUENESS_VIOLATION_PGCODE
        ):
            raise error_cls(**kwargs)
        raise
