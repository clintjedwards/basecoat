from flask import render_template, jsonify
import basecoat.db_utils as db

from app import app


@app.route('/')
def index():
    return render_template('index.html')


@app.route('/api/formulas')
def get_formulas():
    formula_table = db.get_table('Formula')
    formula_list = [formula_table.to_dict() for formula_table in formula_table]

    return jsonify(formula_list)
