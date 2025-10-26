import grpc

from app.proto.user import user_pb2_grpc, user_pb2


class UserClient:
    def __init__(self, address: str):
        self.address = address

    async def check_user_exists(self, user_id: int) -> bool:
        async with grpc.aio.insecure_channel(self.address) as channel:
            stub = user_pb2_grpc.UserV1Stub(channel)
            request = user_pb2.GetRequest(id=user_id)
            try:
                await stub.Get(request, timeout=1.0)
                return True
            except grpc.aio.AioRpcError as e:
                if e.code() == grpc.StatusCode.NOT_FOUND:
                    return False
                raise
