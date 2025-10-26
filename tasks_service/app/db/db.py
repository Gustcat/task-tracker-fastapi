from typing import Annotated

from fastapi import Depends
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker

from app.settings import REAL_DATABASE_URL, settings

engine = create_async_engine(REAL_DATABASE_URL, future=settings.db_echo, echo=True)

async_session = async_sessionmaker(
    engine, expire_on_commit=False, class_=AsyncSession
)


async def get_async_session():
    async with async_session() as session:
        yield session


SessionDep = Annotated[AsyncSession, Depends(get_async_session)]