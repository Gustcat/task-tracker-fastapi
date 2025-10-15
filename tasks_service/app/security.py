from jose import jwt, JWTError
from fastapi import Depends, status, HTTPException
from fastapi.security import OAuth2PasswordBearer

from app.schemas.users import User
from app.settings import settings


token_without_prefix = OAuth2PasswordBearer(tokenUrl="auth/token")  # фейковый url для swagger


def decode_token(token: str):
    try:
        payload = jwt.decode(token, settings.secret_key, algorithms=[settings.algorithm])
        return payload
    except JWTError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid token",
        )


async def get_current_user(token: str = Depends(token_without_prefix)) -> User:
    payload = decode_token(token)
    return User(**payload)

