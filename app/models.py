from app import db

class Formula(db.Model):
    __tablename__ = "Formula"
    id = db.Column(db.Integer, primary_key=True)
    data = db.Column(db.Binary)
