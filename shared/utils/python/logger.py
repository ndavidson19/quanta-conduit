import os
import logging
from logging.handlers import RotatingFileHandler

def setup_logger(name=__name__, log_file='app.log', level=logging.INFO):
    logger = logging.getLogger(name)
    
    # Set log level from environment variable if provided
    log_level = os.getenv('LOG_LEVEL', 'INFO').upper()
    logger.setLevel(getattr(logging, log_level, logging.INFO))

    formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')

    # Console Handler
    console_handler = logging.StreamHandler()
    console_handler.setFormatter(formatter)
    logger.addHandler(console_handler)

    # File Handler
    file_handler = RotatingFileHandler(log_file, maxBytes=10*1024*1024, backupCount=5)
    file_handler.setFormatter(formatter)
    logger.addHandler(file_handler)

    return logger

# Usage
# logger = setup_logger()
# logger.info("This is an info message")