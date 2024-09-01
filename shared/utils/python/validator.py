from pydantic import BaseModel, ValidationError

def validate_model(model: BaseModel, data: dict):
    try:
        return model(**data)
    except ValidationError as e:
        raise BadRequestError(str(e))