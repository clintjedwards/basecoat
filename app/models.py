import datetime

from app import db

#Disable warning: Too few public methods pylint: disable=R0903

class Formula(db.Model):
    __tablename__ = "Formulas"
    id = db.Column(db.Integer, primary_key=True)
    formula_name = db.Column(db.String(64))
    formula_number = db.Column(db.String(64))
    customer_name = db.Column(db.String(64))
    summary = db.Column(db.String(64))
    notes = db.Column(db.Text())
    created_date = db.Column(db.Date, default=datetime.datetime.utcnow)
    last_modified = db.Column(db.Date, default=datetime.datetime.utcnow, onupdate=datetime.datetime.utcnow)


class Colorant(db.Model):
    __tablename__ = "Colorants"
    id = db.Column(db.Integer, primary_key=True)
    formula_id = db.Column(db.Integer, db.ForeignKey("Formulas.id"), nullable=False)
    colorant_name = db.Column(db.String(64))
    amount = db.Column(db.Integer)


class Base(db.Model):
    __tablename__ = "Bases"
    id = db.Column(db.Integer, primary_key=True)
    formula_id = db.Column(db.Integer, db.ForeignKey("Formulas.id"), nullable=False)
    base_name = db.Column(db.String(64))
    product_name = db.Column(db.String(64))
