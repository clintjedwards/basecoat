import datetime

from app import db

#Disable warning: Too few public methods pylint: disable=R0903

class Formula(db.Model):
    __tablename__ = "Formulas"
    id = db.Column(db.Integer, primary_key=True)
    color_name = db.Column(db.String(64))
    color_number = db.Column(db.String(64))
    customer_name = db.Column(db.String(64))
    summary = db.Column(db.String(64))
    notes = db.Column(db.Text())
    created_date = db.Column(db.Date, default=datetime.datetime.utcnow)
    last_modified = db.Column(db.Date, onupdate=datetime.datetime.utcnow)

    def to_dict(self):
        return {
            'id': self.id,
            'color_name': self.color_name,
            'color_number': self.color_number,
            'customer_name': self.customer_name,
            'summary': self.summary,
            'notes': self.notes,
            'created_date': self.created_date,
            'last_modified': self.last_modified,
        }


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
