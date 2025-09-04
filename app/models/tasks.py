from datetime import datetime, date
from enum import StrEnum

from sqlalchemy import String, text, DateTime, Boolean, CheckConstraint, Table, Column, ForeignKey, func
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship


class Base(DeclarativeBase):
    pass


class TimestampModel:
    created_at: Mapped[datetime] = mapped_column(
        DateTime,
        default=datetime.utcnow,
        server_default=func.now(),
        nullable=False
    )
    updated_at: Mapped[datetime] = mapped_column(
        DateTime,
        default=datetime.utcnow,
        onupdate=datetime.utcnow,
        server_default=func.now(),
        server_onupdate=func.now(),
        nullable=False
    )


class TaskStatus(StrEnum):
    NEW = "new"
    TODO = "todo"
    IN_PROGRESS = "in_progress"
    DONE = "done"


statuses = [f"'{status}'" for status in TaskStatus]
check_expr = f"status IN ({', '.join(statuses)})"


class TaskModel(Base, TimestampModel):
    __tablename__ = "task"

    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String(250), nullable=False, unique=True)
    description: Mapped[str | None] = mapped_column(nullable=True)
    status: Mapped[TaskStatus] = mapped_column(String(20), nullable=False, default=TaskStatus.NEW)
    author: Mapped[int] = mapped_column(nullable=False, index=True)
    author_deleted: Mapped[bool] = mapped_column(Boolean, default=False)
    operator: Mapped[int | None] = mapped_column(nullable=True, index=True)
    operator_deleted: Mapped[bool] = mapped_column(Boolean, default=False)
    target_date: Mapped[date | None] = mapped_column(nullable=True)
    completed_at: Mapped[date | None] = mapped_column(nullable=True)
    watchers: Mapped[list["TaskWatcherModel"]] = relationship(
        "TaskWatcherModel",
        back_populates="tasks",
        cascade="all, delete-orphan"
    )

    __table_args__ = (
        CheckConstraint(check_expr, name="check_status"),
    )


class TaskWatcherModel(Base, TimestampModel):
    __tablename__ = "task_watcher"

    task_id: Mapped[int] = mapped_column(ForeignKey("task.id", ondelete="CASCADE"), primary_key=True)
    user_id: Mapped[int] = mapped_column(primary_key=True, index=True)

    tasks: Mapped["TaskModel"] = relationship("TaskModel", back_populates="watchers")
