from flask_script import Manager
from tools.database import MigrateCommand

from app import app

manager = Manager(app)
manager.add_command('db', MigrateCommand)

if __name__ == '__main__':
    manager.run()
