from typing import Dict, Any

class AppError(Exception):
    def __init__(self, code: int, message: str):
        self.code = code
        self.message = message

    def __str__(self):
        return f"Error {self.code}: {self.message}"

    def to_dict(self) -> Dict[str, Any]:
        return {
            "code": self.code,
            "message": self.message
        }

class NotFoundError(AppError):
    def __init__(self, message: str = "Resource not found"):
        super().__init__(404, message)

class BadRequestError(AppError):
    def __init__(self, message: str = "Bad request"):
        super().__init__(400, message)

class InternalServerError(AppError):
    def __init__(self, message: str = "Internal server error"):
        super().__init__(500, message)