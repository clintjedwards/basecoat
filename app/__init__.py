from flask import Flask
from flask_sqlalchemy import SQLAlchemy
from flask_migrate import Migrate

#Disable warning: Import should be placed at the top of the module pylint: disable=C0413

app = Flask(__name__)
app.config.from_object('config.config')

db = SQLAlchemy(app)
migrate = Migrate(app, db)

from app import views, models
