import datetime
from app import db


class Formula(db.Model):
    __tablename__ = "Formula"
    id = db.Column(db.Integer, primary_key=True)
    color_name = db.Column(db.String(64))
    color_number = db.Column(db.Integer)
    customer_name = db.Column(db.String(64))
    created_date = db.Column(db.Date, default=datetime.datetime.utcnow)
    last_modified = db.Column(db.Date, onupdate=datetime.datetime.utcnow)

class Colorant(db.Model):
    __tablename__ = "Colorant"
    formula_id = db.Column(db.Integer, db.ForeignKey("Formula.id"), nullable=False)
    colorant_name = db.Column(db.String(64))
    colorant_amount = db.Column(db.Integer)
