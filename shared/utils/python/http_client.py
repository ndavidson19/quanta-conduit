import httpx

class HTTPClient:
    def __init__(self, base_url: str = "", timeout: float = 10.0):
        self.client = httpx.AsyncClient(base_url=base_url, timeout=timeout)

    async def get(self, url: str, **kwargs):
        return await self.client.get(url, **kwargs)

    async def post(self, url: str, **kwargs):
        return await self.client.post(url, **kwargs)

    async def put(self, url: str, **kwargs):
        return await self.client.put(url, **kwargs)

    async def delete(self, url: str, **kwargs):
        return await self.client.delete(url, **kwargs)

    async def close(self):
        await self.client.aclose()

