from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker

from app import settings

engine = create_async_engine(settings.REAL_DATABASE_URL, future=True, echo=True)

async_session = async_sessionmaker(
    engine, expire_on_commit=False, class_=AsyncSession
)


async def get_async_session():
    async with async_session() as session:
        yield session
