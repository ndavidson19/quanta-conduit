import os
from typing import Any
import yaml

class Config:
    _instance = None

    def __new__(cls):
        if cls._instance is None:
            cls._instance = super(Config, cls).__new__(cls)
            cls._instance._load_config()
        return cls._instance

    def _load_config(self):
        env = os.getenv('ENV', 'development')
        config_path = f"config/{env}.yaml"
        
        with open(config_path, 'r') as config_file:
            self._config = yaml.safe_load(config_file)

    def get(self, key: str) -> Any:
        return self._config.get(key)

config = Config()

# Usage:
# from shared.utils.python.config import config
# db_url = config.get('database_url')