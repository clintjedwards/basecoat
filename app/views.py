from flask import render_template
from app.basecoat import db_utils as db

from app import app


@app.route('/')
def index():
    formula_table = db.get_table('Formula')
    formula_list = [formula_table.to_dict() for formula_table in formula_table]
    return render_template('index.html',
                           formula_list=formula_list)
