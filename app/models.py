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
    colorants = db.Column(db.Text())
    bases = db.Column(db.Text())
    created_date = db.Column(db.Date, default=datetime.datetime.utcnow)
    last_modified = db.Column(db.Date, default=datetime.datetime.utcnow, onupdate=datetime.datetime.utcnow)
