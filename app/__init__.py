from flask import Flask
from flask_sqlalchemy import SQLAlchemy


app = Flask(__name__)
app.config.from_object('config.config')

db = SQLAlchemy(app)

from app import views, models

db.create_all()
